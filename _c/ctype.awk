#!/usr/bin/awk -f
#===============================================================================
#
#          File:  ctype.awk
#
#   Description:
#
#   VIM Version:  7.0+
#        Author:  Solomon Ng,
#  Organization:
#       Version:  1.0
#       Created:  03/23/2018 13:59
#      Revision:  ---
#       License:
#===============================================================================

BEGIN {
    struct_name = ""
    field_name = ""
    type = ""
}

$1 == "}" {
    struct_name = ""
    field_name = ""
    type = ""
}

$2 == "struct" && $3 == "{" {
    struct_name = $1
    field_name = ""
    type = ""
}

struct_name != "" && NF == 2 {
    field_name = $1
    type = $2
}

struct_name == "FuseInHeader" && field_name == "Opcode" {
    type = "OpcodeType"
}

struct_name == "FuseOpenOut" && field_name ==  "Flags" {
    type = "OpenOutFlagType"
}

struct_name == "FuseAttr" && field_name == "Mode" {
    type = "FileModeType"
}

struct_name == "FuseAttr" && field_name == "Rdev" {
    type = "FileModeType"
}

struct_name == "FuseDirent" && field_name == "Type" {
    type = "DirentType"
}

{
    if (field_name != "" && type != "") {
        print field_name, type
    } else {
        print $0
    }
}
