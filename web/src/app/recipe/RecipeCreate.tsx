import "./css/recipe-create.css"
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { api } from "../api";
import { FormText } from "../component/form";
import { TextButton } from "../component/input";
import { Theme } from "../theme";


export default function RecipeCreate() {

  const nav = useNavigate()
  const [name, setName] = useState("")

  return (<div className="recipe-create-root theme-background">
    <div className="recipe-create-workspace theme-focus">
      <div className="recipe-create-header format-font-large">Create a new recipe</div>
      <FormText label="Recipe Name" className="recipe-create-name" value={name}
        inputAttrs={{ autoCapitalize: "words" }} onChange={setName} />
      <TextButton label="Create!" theme={Theme.Primary} className="recipe-create-submit" onClick={e => {
        api.post("recipe/create", { "name": name }).then(rsp => { nav(`/recipe/edit/${rsp.data.uid}`) })
      }} />
    </div>
  </div>)
}