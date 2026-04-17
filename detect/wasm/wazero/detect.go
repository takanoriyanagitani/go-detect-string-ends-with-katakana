package conv

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"os"

	"github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"
)

var (
	ErrNilMem  error = errors.New("nil memory")
	ErrNilFunc error = errors.New("nil function")

	ErrInvalidResults error = errors.New("unexpected num of results")

	ErrUnableToWrite error = errors.New("unable to write the original")
)

type WasmFn struct{ wa.Function }

func (f WasmFn) Detect(
	ctx context.Context,
	hi, lo uint64,
) (bool, error) {
	results, err := f.Function.Call(ctx, hi, lo)

	if nil != err {
		return false, err
	}

	if 1 != len(results) {
		return false, ErrInvalidResults
	}

	var result uint64 = results[0]
	var decoded int32 = wa.DecodeI32(result)

	return 1 == decoded, nil
}

type WasmMod struct{ wa.Module }

func (m WasmMod) Close(ctx context.Context) error {
	if nil == m.Module {
		return nil
	}
	return m.Module.Close(ctx)
}

func (m WasmMod) GetFunction(name string) (WasmFn, error) {
	var fnc wa.Function = m.Module.ExportedFunction(name)
	if nil == fnc {
		return WasmFn{}, ErrNilFunc
	}
	return WasmFn{Function: fnc}, nil
}

func (m WasmMod) GetDetector64() (WasmFn, error) {
	return m.GetFunction("ends_with3zkana64i")
}

type Compiled struct{ wazero.CompiledModule }

func (c Compiled) Close(ctx context.Context) error {
	if nil == c.CompiledModule {
		return nil
	}
	return c.CompiledModule.Close(ctx)
}

type WasmRuntime struct{ wazero.Runtime }

func (r WasmRuntime) Close(ctx context.Context) error {
	if nil == r.Runtime {
		return nil
	}
	return r.Runtime.Close(ctx)
}

func (r WasmRuntime) Compile(
	ctx context.Context,
	wasm []byte,
) (Compiled, error) {
	cmod, err := r.Runtime.CompileModule(ctx, wasm)
	return Compiled{CompiledModule: cmod}, err
}

func (r WasmRuntime) Instantiate(
	ctx context.Context,
	compiled Compiled,
	cfg wazero.ModuleConfig,
) (WasmMod, error) {
	amod, err := r.Runtime.InstantiateModule(
		ctx,
		compiled.CompiledModule,
		cfg,
	)

	return WasmMod{Module: amod}, err
}

type WasmConfig struct{ wazero.RuntimeConfig }

type Detector struct {
	WasmRuntime
	Compiled
	WasmMod
	WasmFn
}

func (c Detector) Close(ctx context.Context) error {
	return errors.Join(
		c.WasmMod.Close(ctx),
		c.Compiled.Close(ctx),
		c.WasmRuntime.Close(ctx),
	)
}

func (c Detector) StdinToStdout(ctx context.Context) error {
	var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	var bwtr *bufio.Writer = bufio.NewWriter(os.Stdout)

	var buf [16]byte
	for scanner.Scan() {
		var line []byte = scanner.Bytes()
		var size int = len(line)

		if size < 9 {
			continue
		}

		var start byte = line[size-9]
		if 0xe3 != start {
			continue
		}

		clear(buf[:])
		copy(buf[:], line[size-9:])

		var first uint64 = binary.LittleEndian.Uint64(buf[:8])
		var last uint64 = binary.LittleEndian.Uint64(buf[8:])

		found, err := c.WasmFn.Detect(ctx, first, last)
		if nil != err {
			return err
		}

		if found {
			_, werr := bwtr.Write(line)
			if nil != werr {
				return werr
			}

			_, werr = bwtr.WriteString("\n")
			if nil != werr {
				return werr
			}
		}
	}

	return bwtr.Flush()
}

type WasmBytes []byte

func (b WasmBytes) ToDetector(
	ctx context.Context,
	rcfg wazero.RuntimeConfig,
	mcfg wazero.ModuleConfig,
) (Detector, error) {
	var rtm wazero.Runtime = wazero.NewRuntimeWithConfig(
		ctx,
		rcfg,
	)
	var detector Detector
	detector.WasmRuntime = WasmRuntime{Runtime: rtm}

	compiled, err := rtm.CompileModule(ctx, b)
	if nil != err {
		return detector, err
	}
	detector.Compiled = Compiled{CompiledModule: compiled}

	instance, err := rtm.InstantiateModule(
		ctx,
		detector.Compiled.CompiledModule,
		mcfg,
	)
	if nil != err {
		return detector, err
	}
	detector.WasmMod = WasmMod{Module: instance}

	det64, err := detector.WasmMod.GetDetector64()
	if nil != err {
		return detector, err
	}
	detector.WasmFn = det64

	return detector, nil
}
