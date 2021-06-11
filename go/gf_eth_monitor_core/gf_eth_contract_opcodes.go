/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_eth_monitor_core

import (
	"strings"
	"encoding/hex"
	"github.com/gloflow/gloflow/go/gf_core"
	eth_asm "github.com/ethereum/go-ethereum/core/asm"
)

/*
0x - Stop and Arithmetic Operations

0x00    STOP        Halts execution
0x01    ADD         Addition operation
0x02    MUL         Multiplication operation
0x03    SUB         Subtraction operation
0x04    DIV         Integer division operation
0x05    SDIV        Signed integer
0x06    MOD         Modulo
0x07    SMOD        Signed modulo
0x08    ADDMOD      Modulo
0x09    MULMOD      Modulo
0x0a    EXP         Exponential operation
0x0b    SIGNEXTEND  Extend length of two's complement signed integer
*/

/*
1x - Comparison & Bitwise Logic Operations

0x10    LT      Lesser-than comparison
0x11    GT      Greater-than comparison
0x12    SLT     Signed less-than comparison
0x13    SGT     Signed greater-than comparison
0x14    EQ      Equality  comparison
0x15    ISZERO  Simple not operator
0x16    AND     Bitwise AND operation
0x17    OR      Bitwise OR operation
0x18    XOR     Bitwise XOR operation
0x19    NOT     Bitwise NOT operation
0x1a    BYTE    Retrieve single byte from word
*/

/*
2x - SHA3

0x20    SHA3    Compute Keccak-256 hash
*/

/*
6x & 7x - Push Operations

0x60    PUSH1   Place 1 byte item on stack
0x61    PUSH2   Place 2-byte item on stack
...
0x7f    PUSH32  Place 32-byte (full word) item on stack
*/

/*
8x - Duplication Operations

0x80    DUP1    Duplicate 1st stack item
0x81    DUP2    Duplicate 2nd stack item
...
0x8f    DUP16   Duplicate 16th stack item
*/

/*
9x - Exchange Operations

0x90    SWAP1   Exchange 1st and 2nd stack items
0x91    SWAP2   Exchange 1st and 3rd stack items
...
0x9f    SWAP16  Exchange 1st and 17th stack items
*/

/*
0xff    SELFDESTRUCT    Halt execution and register account for later deletion
*/

//-------------------------------------------------
type GF_eth__opcode struct {
	Op_and_args_str string
	Addr_hex_str    string
}

// FIX!! - give this struct a better name, too similar to GF_eth__opcode
type gf_eth_opcode struct {
	Important_bool bool // if the opcode is important for GF analysis
}

//-------------------------------------------------
func Eth_contract__get_opcodes(p_bytecode_hex_str string,
	p_runtime *GF_runtime) ([]*GF_eth__opcode, *gf_core.Gf_error) {

	// HEX_DECODE
	code_bytes_lst, err := hex.DecodeString(p_bytecode_hex_str)
	if err != nil {
		gf_err := gf_core.Error__create("failed to decode contract bytecode hex string into byte list",
			"decode_hex",
			map[string]interface{}{
				"bytecode_hex_str": p_bytecode_hex_str,
			},
			err, "gf_eth_monitor_core", p_runtime.Runtime_sys)
		return nil, gf_err
	}

	// eth_asm.PrintDisassembled(code_str)

	// DISASSEMBLE
	output_lst, err := eth_asm.Disassemble(code_bytes_lst)
	if err != nil {
		error_defs_map := Error__get_defs()
		gf_err := gf_core.Error__create_with_defs("failed to disassemble contract hex bytecode",
			"eth_contract__disassemble",
			map[string]interface{}{
				"bytecode_hex_str": p_bytecode_hex_str,
			},
			err, "gf_eth_monitor_core", error_defs_map, 1, p_runtime.Runtime_sys)
		return nil, gf_err
	}

	opcodes_lst := []*GF_eth__opcode{}
	for _, opcode_str := range output_lst {

		s_lst := strings.Split(strings.TrimSpace(opcode_str), ": ")
		opcode_addr_hex_str := s_lst[0]
		opcode_and_args_str := s_lst[1]

		// fmt.Printf("%s - %s\n", opcode_addr_hex_str, opcode_and_args_str)

		// GF_OPCODE
		gf_opcode := &GF_eth__opcode{
			Op_and_args_str: opcode_and_args_str,
			Addr_hex_str:    opcode_addr_hex_str,
		}
		opcodes_lst = append(opcodes_lst, gf_opcode)
	}

	return opcodes_lst, nil
}

//-------------------------------------------------
func Eth_opscodes__get_tables() map[string]map[string]gf_eth_opcode {


	/*
	3x - Environmental Information

	0x30    ADDRESS         Get address of currently executing account
	0x31    BALANCE         Get balance of the given account
	0x32    ORIGIN          Get execution origination address
	0x33    CALLER          Get caller address. This is the address of the account that is directly responsible for this execution
	0x34    CALLVALUE       Get deposited value by the instruction/transaction responsible for this execution
	0x35    CALLDATALOAD    Get input data of current environment
	0x36    CALLDATASIZE    Get size of input data in current environment
	0x37    CALLDATACOPY    Copy input data in current environment to memory This pertains to the input data passed with the message call instruction or transaction
	0x38    CODESIZE        Get size of code running in current environment
	0x39    CODECOPY        Copy code running in current environment to memory
	0x3a    GASPRICE        Get price of gas in current environment
	0x3b    EXTCODESIZE     Get size of an account's code
	0x3c    EXTCODECOPY     Copy an account's code to memory
	*/
	ops__env_map := map[string]gf_eth_opcode{
		"ADDRESS":      gf_eth_opcode{Important_bool: false},
		"BALANCE":      gf_eth_opcode{Important_bool: false},
		"ORIGIN":       gf_eth_opcode{Important_bool: true},
		"CALLER":       gf_eth_opcode{Important_bool: true},
		"CALLVALUE":    gf_eth_opcode{Important_bool: true},
		"CALLDATALOAD": gf_eth_opcode{Important_bool: true},
		"CALLDATASIZE": gf_eth_opcode{Important_bool: false},
		"CALLDATACOPY": gf_eth_opcode{Important_bool: false},
		"CODESIZE":     gf_eth_opcode{Important_bool: false},
		"CODECOPY":     gf_eth_opcode{Important_bool: false},
		"GASPRICE":     gf_eth_opcode{Important_bool: false},
		"EXTCODESIZE":  gf_eth_opcode{Important_bool: false},
		"EXTCODECOPY":  gf_eth_opcode{Important_bool: false},
	}

	/*
	4x - Block Information

	0x40    BLOCKHASH   Get the hash of one of the 256 most recent complete blocks
	0x41    COINBASE    Get the block's beneficiary address
	0x42    TIMESTAMP   Get the block's timestamp
	0x43    NUMBER      Get the block's number
	0x44    DIFFICULTY  Get the block's difficulty
	0x45    GASLIMIT    Get the block's gas limit
	*/
	/*ops__block_map := map[string]interface{}{
		"BLOCKHASH":  gf_eth_opcode{Important_bool: false},
		"COINBASE":   gf_eth_opcode{Important_bool: false},
		"TIMESTAMP":  gf_eth_opcode{Important_bool: false},
		"NUMBER":     gf_eth_opcode{Important_bool: false},
		"DIFFICULTY": gf_eth_opcode{Important_bool: false},
		"GASLIMIT":   gf_eth_opcode{Important_bool: false},
	}*/

	/*
	5x - Stack, Memory, Storage and Flow Operations

	0x50    POP         Remove item from stack
	0x51    MLOAD       Load word from memory
	0x52    MSTORE      Save word to memory
	0x53    MSTORE8     Save byte to memory
	0x54    SLOAD       Load word from storage
	0x55    SSTORE      Save word to storage
	0x56    JUMP        Alter the program counter
	0x57    JUMPI       Conditionally alter the program counter
	0x58    PC          Get the value of the program counter prior to the increment
	0x59    MSIZE       Get the size of active memory in bytes
	0x5a    GAS         Get the amount of available gas, including the corresponding reduction
	0x5b    JUMPDEST    Mark a valid destination for jumps
	*/
	ops__memory_map := map[string]gf_eth_opcode{
		"POP":      gf_eth_opcode{Important_bool: false},
		"MLOAD":    gf_eth_opcode{Important_bool: true},
		"MSTORE":   gf_eth_opcode{Important_bool: true},
		"MSTORE8":  gf_eth_opcode{Important_bool: true},
		"SLOAD":    gf_eth_opcode{Important_bool: true},
		"SSTORE":   gf_eth_opcode{Important_bool: true},
		"JUMP":     gf_eth_opcode{Important_bool: false},
		"JUMPI":    gf_eth_opcode{Important_bool: false},
		"PC":       gf_eth_opcode{Important_bool: false},
		"MSIZE":    gf_eth_opcode{Important_bool: false},
		"GAS":      gf_eth_opcode{Important_bool: false},
		"JUMPDEST": gf_eth_opcode{Important_bool: false},
	}

	/*
	fx - System operations

	0xf0    CREATE          Create a new account with associated code
	0xf1    CALL            Message-call into an account
	0xf2    CALLCODE        Message-call into this account with alternative account's code
	0xf3    RETURN          Halt execution returning output data
	0xf4    DELEGATECALL    Message-call into this account with an alternative account's code, but persisting the current values for `sender` and `value`
	*/
	ops__sys_map := map[string]gf_eth_opcode{
		"CREATE":       gf_eth_opcode{Important_bool: true},
		"CALL":         gf_eth_opcode{Important_bool: true},
		"CALLCODE":     gf_eth_opcode{Important_bool: true},
		"RETURN":       gf_eth_opcode{Important_bool: true},
		"DELEGATECALL": gf_eth_opcode{Important_bool: true},
	}


	/*
	ax - Logging Operations

	0xa0    LOG0    Append log record with no topics
	0xa1    LOG1    Append log record with one topic
	...
	0xa4    LOG4    Append log record with four topics
	*/
	ops__log_map := map[string]gf_eth_opcode{
		"LOG0": gf_eth_opcode{Important_bool: true},
		"LOG1": gf_eth_opcode{Important_bool: true},
		"LOG2": gf_eth_opcode{Important_bool: true},
		"LOG3": gf_eth_opcode{Important_bool: true},
		"LOG4": gf_eth_opcode{Important_bool: true},
	}





	ops_map := map[string]map[string]gf_eth_opcode{
		"env_map":    ops__env_map,
		"memory_map": ops__memory_map,
		"sys_map":    ops__sys_map,
		"log_map":    ops__log_map,
	}





	return ops_map

}