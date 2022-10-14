export const Units = ["cup", "tsp", "Tbsp", "lbs", "kg", "g", "mL", "L", "pt", "pt", "gal", ""]

export type RecipeModel = {
  uid: string,
  created_at: string,
  updated_at: string,
  is_trashed: boolean,
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
  updated_at: string,
  meal: string,
  order: number,
  unit: string | null,
  quantity: string | number,
  item: string,
}

export type InstructionModel = {
  uid: string,
  user_uid: string,
  recipe_uid: string,
  created_at: string,
  updated_at: string,
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
    is_trashed: model.is_trashed,
    description: model.description,
    meal: model.meal,
    ingredients: model.ingredients.map(i => { return { quantity: i.quantity, unit: i.unit, item: i.item } }),
    instructions: model.instructions.map(i => i.text),
  }
}

