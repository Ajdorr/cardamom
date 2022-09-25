import { useEffect, useState } from "react"
import { api } from "../api"
import { ImageButton } from "../component/input"
import { RecipeModel } from "./schema"

function RecipeTrashList() {
  const [recipes, setRecipes] = useState<RecipeModel[]>([])

  useEffect(() => {
    api.post("recipe/trash").then(rsp => {
      setRecipes(rsp.data)
    }).catch(e => console.log(e))
  }, [])

  return (<div className="recipe-trash-list-root">

    <div className="recipe-trash-list-recipes">
      {recipes.length > 0 ? recipes.map(r => {
        return (
          <div key={r.uid} className="recipe-trash-list-element-root">
            <div className="recipe-trash-list-element-name"><span>{r.name}</span></div>
            <ImageButton alt="trash recipe" src="/icons/undelete.svg" className="recipe-list-element-untrash"
              onClick={e => {
                api.post("recipe/update", { uid: r.uid, is_trashed: false }).then(rsp => {
                  let newRecipes = [...recipes].filter(nr => { return nr.uid !== r.uid })
                  setRecipes(newRecipes)
                }).catch(console.log)
              }}
            />
          </div>
        )
      }) : <div className="recipe-list-empty">Trash is empty</div>}
    </div>
  </div>)
}

export default RecipeTrashList