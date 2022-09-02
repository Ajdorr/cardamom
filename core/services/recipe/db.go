package recipe

import gonanoid "github.com/matoous/go-nanoid/v2"

func generateRecipeUid() string {
	return gonanoid.Must()
}

func generateInstrUid() string {
	return gonanoid.Must(24)
}

func generateIngreUid() string {
	return gonanoid.Must(24)
}
