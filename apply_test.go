package dragon

import (
	"testing"
)

// Test that two different sequences of moves involving en passant but leading to the same board have the same result
func TestHashEpApplication(t *testing.T) {
	b1 := ParseFen(Startpos)
	b1.Apply(parseMove("e2e4"))
	b1.Apply(parseMove("b8c6"))
	b2 := ParseFen("r1bqkbnr/pppppppp/2n5/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
	if b1.Hash() != b2.Hash() {
		t.Error("Fen changed.")
	}
}

func TestApplyUnapply(t *testing.T) {
	movesMap := map[string]Move{
		// ordinary move
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0": parseMove("e2e4"),
		// castle 1: white short
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 0": parseMove("e1g1"),
		// castle 2: black long, without kingside rights
		"r3kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R b KQq - 0 0": parseMove("e8c8"),
		// rook move strips castling rights
		"r3k2r/1pppppp1/8/8/8/8/1PPPPPPP/R3K2R w KQq - 0 0": parseMove("a1b1"),
		// if no castling rights, rook move has no effect
		"r3k2r/1pppppp1/8/8/8/8/1PPPPPPP/R3K2R b KQq - 0 0": parseMove("h8h7"),
		// en passant capture
		"r3k3/1ppp1ppr/8/3Pp3/8/8/1PP1PPPP/R3K2R w - e6 3 0": parseMove("d5e6"),
		"r3k3/1ppp1ppr/8/8/2Pp4/8/1P2PPPP/R3K2R b - c3 0 0":  parseMove("d4c3"),
		// pawn push updates en passant: white
		"2kr1bnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQ - 0 0": parseMove("a2a4"),
		// pawn push updates en passant: black
		// promotion 1: white to queen
		"r3k3/1pp3P1/4N3/3b4/8/2p5/1P2PP1P/R3K2R w - - 0 0": parseMove("g7g8q"),
		// promotion 2: black to knight
		"r3k1Q1/1pp5/4N3/3b4/8/2p5/1P2PP1p/R3K3 b - - 0 0": parseMove("h2h1n"),
		// promotion-capture: black underpromotion
		"r3k1Q1/1pp5/4N3/3br3/8/2p3n1/1p2PP2/R1B1K2n b - - 0 0": parseMove("b2c1b"),
		// capture: black king captures white knight
		"r3k1Q1/1pp2p2/4Nk2/3br3/8/2p3n1/4PP2/R1b1K2n b - - 0 0": parseMove("f6e6"),
		// king: strip castle rights bug
		"rnbqkbnr/ppp1pppp/8/3p4/8/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2": parseMove("e1d2"),
		// king: e.p. bug
		"rnbqkbnr/ppp1pppp/8/3p4/8/8/PPP1PPPP/RNBQKBNR w KQkq d6 0 2": parseMove("e1d2"),
		// castling in kiwipete
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0": parseMove("e1c1"),
		// a buggy position
		"r3k3/p1ppqpb1/bn2pnpr/3PN3/1p2P3/5Q1p/PPPBBPPP/RN2K2R w KQq - 0 0": parseMove("d2h6"),
		// castling bug
		"r3k2r/p1ppqpb1/1n2pnp1/1b1PN3/1p2P3/P1N2Q1p/1PPBBPPP/R3K2R w KQkq - 0 0": parseMove("e1g1"),
		// ucapture of an underpromoted rook shouldn't strip castle rights
		"r3k2r/Pppp1ppp/1b3nbN/nPB5/B1P1P3/q4N2/P2P2PP/r2Q1RK1 w kq - 0 0": parseMove("d1a1"),
		// Promoted rook captures pawn, shouldn't strip castle rights
		"r3k2r/Pppp1ppp/1b3nbN/nPB5/2P1P3/qB3N2/P2P2PP/r2Q1RK1 b kq - 0 0": parseMove("a1a2"),
		// Moving toward a forced mate
		"5k2/5p2/5P2/8/8/2r5/2rR2K1/4B2R w - - 0 1": parseMove("h1h8"),
	}
	results := map[string]string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0":                "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 0",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQkq - 0 0":                 "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQ1RK1 b kq - 1 0",
		"r3kbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R b KQq - 0 0":                    "2kr1bnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQ - 1 1",
		"r3k2r/1pppppp1/8/8/8/8/1PPPPPPP/R3K2R w KQq - 0 0":                       "r3k2r/1pppppp1/8/8/8/8/1PPPPPPP/1R2K2R b Kq - 1 0",
		"r3k2r/1pppppp1/8/8/8/8/1PPPPPPP/R3K2R b KQq - 0 0":                       "r3k3/1ppppppr/8/8/8/8/1PPPPPPP/R3K2R w KQq - 1 1",
		"r3k3/1ppp1ppr/8/3Pp3/8/8/1PP1PPPP/R3K2R w - e6 3 0":                      "r3k3/1ppp1ppr/4P3/8/8/8/1PP1PPPP/R3K2R b - - 0 0",
		"r3k3/1ppp1ppr/8/8/2Pp4/8/1P2PPPP/R3K2R b - c3 0 0":                       "r3k3/1ppp1ppr/8/8/8/2p5/1P2PPPP/R3K2R w - - 0 1",
		"r3k3/1pp3P1/4N3/3b4/8/2p5/1P2PP1P/R3K2R w - - 0 0":                       "r3k1Q1/1pp5/4N3/3b4/8/2p5/1P2PP1P/R3K2R b - - 0 0",
		"r3k1Q1/1pp5/4N3/3b4/8/2p5/1P2PP1p/R3K3 b - - 0 0":                        "r3k1Q1/1pp5/4N3/3b4/8/2p5/1P2PP2/R3K2n w - - 0 1",
		"r3k1Q1/1pp5/4N3/3br3/8/2p3n1/1p2PP2/R1B1K2n b - - 0 0":                   "r3k1Q1/1pp5/4N3/3br3/8/2p3n1/4PP2/R1b1K2n w - - 0 1",
		"r3k1Q1/1pp2p2/4Nk2/3br3/8/2p3n1/4PP2/R1b1K2n b - - 0 0":                  "r3k1Q1/1pp2p2/4k3/3br3/8/2p3n1/4PP2/R1b1K2n w - - 0 1",
		"2kr1bnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK2R w KQ - 0 0":                    "2kr1bnr/pppppppp/8/8/P7/8/1PPPPPPP/RNBQK2R b KQ a3 0 0",
		"rnbqkbnr/ppp1pppp/8/3p4/8/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2":              "rnbqkbnr/ppp1pppp/8/3p4/8/8/PPPKPPPP/RNBQ1BNR b kq - 1 2",
		"rnbqkbnr/ppp1pppp/8/3p4/8/8/PPP1PPPP/RNBQKBNR w KQkq d6 0 2":             "rnbqkbnr/ppp1pppp/8/3p4/8/8/PPPKPPPP/RNBQ1BNR b kq - 1 2",
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 0":    "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/2KR3R b kq - 1 0",
		"r3k3/p1ppqpb1/bn2pnpr/3PN3/1p2P3/5Q1p/PPPBBPPP/RN2K2R w KQq - 0 0":       "r3k3/p1ppqpb1/bn2pnpB/3PN3/1p2P3/5Q1p/PPP1BPPP/RN2K2R b KQq - 0 0",
		"r3k2r/p1ppqpb1/1n2pnp1/1b1PN3/1p2P3/P1N2Q1p/1PPBBPPP/R3K2R w KQkq - 0 0": "r3k2r/p1ppqpb1/1n2pnp1/1b1PN3/1p2P3/P1N2Q1p/1PPBBPPP/R4RK1 b kq - 1 0",
		"r3k2r/Pppp1ppp/1b3nbN/nPB5/B1P1P3/q4N2/P2P2PP/r2Q1RK1 w kq - 0 0":        "r3k2r/Pppp1ppp/1b3nbN/nPB5/B1P1P3/q4N2/P2P2PP/Q4RK1 b kq - 0 0",
		"r3k2r/Pppp1ppp/1b3nbN/nPB5/2P1P3/qB3N2/P2P2PP/r2Q1RK1 b kq - 0 0":        "r3k2r/Pppp1ppp/1b3nbN/nPB5/2P1P3/qB3N2/r2P2PP/3Q1RK1 w kq - 0 1",
		"5k2/5p2/5P2/8/8/2r5/2rR2K1/4B2R w - - 0 1":                               "5k1R/5p2/5P2/8/8/2r5/2rR2K1/4B3 b - - 1 1",
	}
	for k, v := range movesMap {
		b := ParseFen(k)
		oldHash := b.Hash()
		fenBefore := b.ToFen()
		fenAfter := b.ToFen()
		if fenBefore != k {
			t.Error("Fen changed during parsing for board", k)
		}
		if fenBefore != fenAfter {
			t.Error("Fen changed during generation for board", k)
		}
		unapply := b.Apply(v)
		if b.ToFen() != results[k] {
			t.Error("Move application of\n", &v, "\ndidn't produce expected result for\n", k, "->\n",
				results[k], "\nInstead, we got:\n", b.ToFen())
		}
		if b.Hash() != recomputeBoardHash(&b) {
			t.Error("Move apply changed board hash from expected result",
				"\nwith move", &v)
		}
		unapply()
		newHash := b.Hash()
		if oldHash != newHash {
			t.Error("(0) Move unapply (or previous apply) changed board hash for:\n",
				b.ToFen(), "\nwith move", &v)
		}
		if newHash != recomputeBoardHash(&b) {
			t.Error("(1) Move unapply (or previous apply) changed board hash for:\n",
				b.ToFen(), "\nwith move", &v)
		}
		if k != b.ToFen() {
			t.Error("Board changed during unapply for\n", k, "\nResult was\n", b.ToFen(),
				"\nwith move", &v)
		}
		/*movesList := b.GenerateLegalMoves()
		for _, mv := range movesList {
			oldHash := b.Hash()
			unapply := b.Apply(mv)
			if b.Hash() != recomputeBoardHash(&b) {
				t.Error("(1) Move apply changed board hash from expected result")
			}
			unapply()
			newHash := b.Hash()
			if b.ToFen() != k {
				t.Error("Move unapply (or previous apply) changed board for:\n",
					b.ToFen(), "\nand move", &mv)
			}
			if oldHash != newHash {
				t.Error("(2) Move unapply (or previous apply) changed board hash for:\n",
					b.ToFen(), "\nand move", &mv)
			}
			if newHash != recomputeBoardHash(&b) {
				t.Error("(3) Move unapply (or previous apply) changed board hash for:\n",
					b.ToFen(), "\nand move", &mv)
			}
		}*/
	}
}
