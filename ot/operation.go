package ot

import (
	"CollabDoc-go/model/database"
	"errors"
	"log"
	"unicode/utf8"
	// 替换成实际路径
)

func ApplyAndTransform(op database.Operation, baseText string, opsHistory []database.Operation) (string, database.Operation, error) {
	transformed := op

	// 逐个变形
	for _, histOp := range opsHistory {
		transformed = Transform(transformed, histOp)
	}

	// 应用到 baseText 上
	newText, err := Apply(transformed, baseText)
	if err != nil {
		return "", database.Operation{}, err
	}

	return newText, transformed, nil
}

// Apply 操作应用文本，返回修改后文本或错误
func Apply(op database.Operation, text string) (string, error) {
	runes := []rune(text)
	log.Printf("[OT] Apply: text len=%d, op pos=%d, op type=%s, op text=%q", utf8.RuneCountInString(text), op.Position, op.Type, op.Text)

	switch op.Type {
	case "insert":
		existing := runes[op.Position : op.Position+len([]rune(op.Text))]
		if string(existing) == op.Text {
			return string(runes), nil // 重复插入，跳过
		}
		if op.Position < 0 || op.Position > len(runes) {
			return "", errors.New("invalid insert position")
		}
		ins := []rune(op.Text)
		res := append(runes[:op.Position], append(ins, runes[op.Position:]...)...)
		return string(res), nil

	case "delete":
		start := op.Position
		delRunes := []rune(op.Text)
		delLen := len(delRunes)
		end := start + delLen

		if start < 0 || end > len(runes) || start > end {
			return "", errors.New("invalid delete range")
		}

		// 可选：校验被删除文本一致性，防止数据冲突
		if string(runes[start:end]) != op.Text {
			return "", errors.New("delete text mismatch with original content")
		}

		res := append(runes[:start], runes[end:]...)
		return string(res), nil

	case "sync":
		return op.Text, nil

	default:
		return "", errors.New("unknown op type")
	}
}

// Transform 把 op1 转换到包含 op2 的上下文下
func Transform(op1, op2 database.Operation) database.Operation {
	op1Copy := op1 // 避免直接修改传入值

	if op2.Type == "insert" {
		if op1Copy.Position > op2.Position ||
			(op1Copy.Position == op2.Position && op1Copy.Type == "insert") {
			op1Copy.Position += len([]rune(op2.Text))
		}
		log.Printf("Transform: op1 pos=%d => %d by op2 pos=%d type=%s",
			op1.Position, op1Copy.Position, op2.Position, op2.Type)
	} else if op2.Type == "delete" {
		delStart := op2.Position
		delLen := len([]rune(op2.Text))
		delEnd := delStart + delLen

		if op1Copy.Position >= delEnd {
			op1Copy.Position -= delLen
		} else if op1Copy.Position >= delStart && op1Copy.Position < delEnd {
			op1Copy.Position = delStart
			if op1Copy.Type == "delete" {
				op1Copy.Text = ""
			}
		}
	}

	return op1Copy
}
