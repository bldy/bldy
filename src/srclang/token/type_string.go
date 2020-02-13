// Code generated by "stringer -type Type"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EOF-0]
	_ = x[NEWLINE-1]
	_ = x[ERROR-2]
	_ = x[ADD-3]
	_ = x[SUB-4]
	_ = x[MUL-5]
	_ = x[QUO-6]
	_ = x[REM-7]
	_ = x[AND-8]
	_ = x[OR-9]
	_ = x[XOR-10]
	_ = x[SHL-11]
	_ = x[SHR-12]
	_ = x[AND_NOT-13]
	_ = x[ADD_ASSIGN-14]
	_ = x[SUB_ASSIGN-15]
	_ = x[MUL_ASSIGN-16]
	_ = x[QUO_ASSIGN-17]
	_ = x[REM_ASSIGN-18]
	_ = x[AND_ASSIGN-19]
	_ = x[OR_ASSIGN-20]
	_ = x[XOR_ASSIGN-21]
	_ = x[SHL_ASSIGN-22]
	_ = x[SHR_ASSIGN-23]
	_ = x[AND_NOT_ASSIGN-24]
	_ = x[LAND-25]
	_ = x[LOR-26]
	_ = x[ARROW-27]
	_ = x[INC-28]
	_ = x[DEC-29]
	_ = x[EQL-30]
	_ = x[LSS-31]
	_ = x[GTR-32]
	_ = x[ASSIGN-33]
	_ = x[NOT-34]
	_ = x[NEQ-35]
	_ = x[LEQ-36]
	_ = x[GEQ-37]
	_ = x[DEFINE-38]
	_ = x[ELLIPSIS-39]
	_ = x[LPAREN-40]
	_ = x[LBRACK-41]
	_ = x[LBRACE-42]
	_ = x[COMMA-43]
	_ = x[PERIOD-44]
	_ = x[RPAREN-45]
	_ = x[RBRACK-46]
	_ = x[RBRACE-47]
	_ = x[SEMICOLON-48]
	_ = x[COLON-49]
	_ = x[TYPE-50]
	_ = x[MODULE-51]
	_ = x[LET-52]
	_ = x[IDENT-53]
	_ = x[STRING-54]
	_ = x[INT-55]
	_ = x[FLOAT-56]
	_ = x[HEX-57]
	_ = x[FUNC-58]
	_ = x[keywords_end-59]
}

const _Type_name = "EOFNEWLINEERRORADDSUBMULQUOREMANDORXORSHLSHRAND_NOTADD_ASSIGNSUB_ASSIGNMUL_ASSIGNQUO_ASSIGNREM_ASSIGNAND_ASSIGNOR_ASSIGNXOR_ASSIGNSHL_ASSIGNSHR_ASSIGNAND_NOT_ASSIGNLANDLORARROWINCDECEQLLSSGTRASSIGNNOTNEQLEQGEQDEFINEELLIPSISLPARENLBRACKLBRACECOMMAPERIODRPARENRBRACKRBRACESEMICOLONCOLONTYPEMODULELETIDENTSTRINGINTFLOATHEXFUNCkeywords_end"

var _Type_index = [...]uint16{0, 3, 10, 15, 18, 21, 24, 27, 30, 33, 35, 38, 41, 44, 51, 61, 71, 81, 91, 101, 111, 120, 130, 140, 150, 164, 168, 171, 176, 179, 182, 185, 188, 191, 197, 200, 203, 206, 209, 215, 223, 229, 235, 241, 246, 252, 258, 264, 270, 279, 284, 288, 294, 297, 302, 308, 311, 316, 319, 323, 335}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
