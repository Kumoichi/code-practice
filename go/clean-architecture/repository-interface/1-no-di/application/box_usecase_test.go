package application

import (
	"strings"
	"testing"
)

// DIをしていない悪い例。NewBoxUseCase()は引数が0個で、Mockを渡す場所が
// そもそも存在しない。UseCase(ビジネスロジック)のテストのはずなのに、
// 直接本物のDBに繋いでテストするしかない。比較: 3-di-interface/application/box_usecase_test.go

func TestBoxUseCase_IsLarge_5は4以上なのでtrue(t *testing.T) {
	useCase := NewBoxUseCase() // ← 引数なし。DIできないのでMockを渡しようがない

	got, err := useCase.IsLarge(1) // seed data: id=1, number=5
	if err != nil {
		t.Skipf("postgres is required for this test (run `docker compose up -d` first): %v", err)
	}
	if !got {
		t.Errorf("IsLarge() = false, want true")
	}
}

func TestBoxUseCase_IsLarge_2は4未満なのでfalse(t *testing.T) {
	useCase := NewBoxUseCase()

	got, err := useCase.IsLarge(2) // seed data: id=2, number=2
	if err != nil {
		t.Skipf("postgres is required for this test (run `docker compose up -d` first): %v", err)
	}
	if got {
		t.Errorf("IsLarge() = true, want false")
	}
}

func TestBoxUseCase_IsLarge_存在しないIDはエラー(t *testing.T) {
	useCase := NewBoxUseCase()

	_, err := useCase.IsLarge(9999) // does not exist
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	// A connection failure also returns a non-nil err, so "err == nil" alone
	// cannot tell "not found" apart from "postgres is unreachable". Without
	// this check, this test would falsely PASS with postgres stopped.
	if !strings.Contains(err.Error(), "not found") {
		t.Skipf("postgres is required for this test (run `docker compose up -d` first): %v", err)
	}
}
