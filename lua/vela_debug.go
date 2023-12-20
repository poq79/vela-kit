package lua

type Sample struct {
	FnIsG             bool
	FnSourceName      string
	FnLineDefined     int
	FnLastLineDefined int

	//parent
	ParentIdx        int
	ParentPc         int
	ParentBase       int
	ParentLocalBase  int
	ParentReturnBase int
	ParentNArgs      int
	ParentNRet       int
	ParentTailCall   int
}

func (dbg *Debug) Sample() *Sample {
	if dbg.frame == nil {
		return nil
	}

	return &Sample{
		FnIsG:             dbg.frame.Fn.IsG,
		FnSourceName:      dbg.frame.Fn.Proto.SourceName,
		FnLineDefined:     dbg.frame.Fn.Proto.LineDefined,
		FnLastLineDefined: dbg.frame.Fn.Proto.LastLineDefined,
		ParentIdx:         dbg.frame.Parent.Idx,
		ParentPc:          dbg.frame.Parent.Pc,
		ParentBase:        dbg.frame.Parent.Base,
		ParentLocalBase:   dbg.frame.Parent.LocalBase,
		ParentReturnBase:  dbg.frame.Parent.ReturnBase,
		ParentNArgs:       dbg.frame.Parent.NArgs,
		ParentNRet:        dbg.frame.Parent.NRet,
		ParentTailCall:    dbg.frame.Parent.TailCall,
	}
}
