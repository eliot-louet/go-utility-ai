package considerations

import (
	"github.com/eliot-louet/go-utility-ai/ai"
	"github.com/eliot-louet/go-utility-ai/ai/curves"
)

const BoredomConsiderationID ai.ConsiderationID = "OwnBoredom"
const InverseBoredomConsiderationID ai.ConsiderationID = "OwnInverseBoredom"

func OwnBoredom() *ai.Consideration {
	return &ai.Consideration{
		ID: BoredomConsiderationID,
		InputFunc: func(ctx *ai.Context, _ ai.Target) float64 {
			var boredom, ok = ctx.Self.GetValue("boredom")

			if !ok {
				return 0
			}

			return float64(boredom.(int)) / 100
		},

		MinValue: 0,
		MaxValue: 1,

		ResponseCurve: curves.Linear{
			M: 1,
			B: 0,
		},
	}
}

func OwnInverseBoredom() *ai.Consideration {
	return &ai.Consideration{
		ID: InverseBoredomConsiderationID,
		InputFunc: func(ctx *ai.Context, _ ai.Target) float64 {
			var boredom, ok = ctx.Self.GetValue("boredom")

			if !ok {
				return 0
			}
			return float64(boredom.(int)) / 100
		},

		MinValue: 0,
		MaxValue: 1,

		ResponseCurve: curves.Linear{
			M: -1,
			B: 1,
		},
	}
}
