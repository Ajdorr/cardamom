from selenium.webdriver.support.select import Select
from selenium.webdriver.common.by import By
from selenium.webdriver.remote.webdriver import WebDriver


def verify(d: WebDriver, info: dict):

  search = f"//div[contains(@class,'recipe-single-name')]/*/input[@value='{info['name']}']"
  assert len(d.find_elements(By.XPATH, search)) == 1

  meal = Select(
      d.find_element(By.CSS_SELECTOR, ".recipe-single-meal select")
  ).first_selected_option
  assert meal.text == info['meal']

  # search = f"//div[contains(@class,'recipe-single-desc')]/textarea[text()='{info['desc']}']"
  ele = d.find_element(
      By.XPATH, "//div[contains(@class,'recipe-single-desc')]/textarea")
  assert ele.text == info['desc']

  quantity = "(//div[contains(@class,'recipe-ingredient-quantity')])[%d]/input[@value='%s']"
  unit = "(//select[contains(@class,'recipe-ingredient-unit')])[%d]"
  item = "(//div[contains(@class,'recipe-ingredient-item')])[%d]/input[@value='%s']"
  optional = "(//div[contains(@class,'recipe-ingredient-optional')])[%d]/input"

  assert len(d.find_elements(By.CSS_SELECTOR,
             ".recipe-ingredient-root")) == len(info["ingre"])

  for i, ingre in enumerate(info["ingre"], 1):
    u = Select(d.find_element(By.XPATH, unit % i)).first_selected_option
    assert u.text == ingre["unit"]
    item_value = ", ".join(
        [ingre["item"], ingre["modifier"]]) if "modifier" in ingre else ingre["item"]
    assert len(d.find_elements(By.XPATH, item % (i, item_value))) == 1
    assert len(d.find_elements(By.XPATH, quantity %
               (i, ingre["quantity"]))) == 1
    is_optional = d.find_element(
        By.XPATH, optional % i).get_attribute("checked") != None
    assert is_optional == ingre.get("optional", False)

  instrs = d.find_elements(By.CSS_SELECTOR, ".recipe-instruction-list li")
  assert [i.text for i in instrs] == info["instr"].split("\n")
