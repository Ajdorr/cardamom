package recipe

import gonanoid "github.com/matoous/go-nanoid/v2"

func generateRecipeUid() string {
	return gonanoid.Must()
}

func generateIngreUid() string {
	return gonanoid.Must(24)
}
