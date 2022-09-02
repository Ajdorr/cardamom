export const Units = ["cup", "tsp", "Tbsp", "kg", "g", "pt", "pt", "gal"]

export type RecipeModel = {
  uid: string,
  created_at: string,
  update_at: string,
  user_uid: string,
  name: string,
  description: string,
  meal: string,
  instructions: InstructionModel[],
  ingredients: IngredientModel[],
}

export type IngredientModel = {
  uid: string,
  user_uid: string,
  recipe_uid: string,
  created_at: string,
  update_at: string,
  meal: string,
  order: number,
  unit: string,
  quantity: string | number,
  item: string,
}

export type InstructionModel = {
  uid: string,
  user_uid: string,
  recipe_uid: string,
  created_at: string,
  update_at: string,
  meal: string,
  order: number,
  text: string,
}

export const MealTypes = new Map<string, string>([
  ["breakfast", "Breakfast"],
  ["lunch", "Lunch"],
  ["dinner", "Dinner"],
  ["dessert", "Dessert"],
])

export function CreateRecipeRequest(model: RecipeModel): any {
  return {
    name: model.name,
    description: model.description,
    meal: model.meal,
    instructions: model.instructions.map(i => i.text),
    ingredients: model.ingredients.map(i => { return { quantity: i.quantity, unit: i.unit, item: i.item } }),
  }
}

export function UpdateRecipeRequest(model: RecipeModel): any {
  return {
    uid: model.uid,
    name: model.name,
    description: model.description,
    meal: model.meal,
    ingredients: model.ingredients.map(i => { return { quantity: i.quantity, unit: i.unit, item: i.item } }),
    instructions: model.instructions.map(i => i.text),
  }
}

