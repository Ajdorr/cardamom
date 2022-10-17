export const InventoryCategories = new Map<string, string>([
  ["cooking", "Cooking"],
  ["spices", "Spices"],
  ["sauces", "Sauces and Condiments"],
  ["non-perishables", "Non-Perishables"],
  ["non-cooking", "Non-Cooking"],
])

export type InventoryItemModel = {
  uid: string
  created_at: string
  updated_at: string
  user_uid: string
  item: string
  in_stock: boolean
  category: string
}

