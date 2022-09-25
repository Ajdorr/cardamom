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

  assert len(d.find_elements(By.CSS_SELECTOR,
             ".recipe-ingredient-root")) == len(info["ingre"])

  for i, ingre in enumerate(info["ingre"], 1):
    u = Select(d.find_element(By.XPATH, unit % i)).first_selected_option
    assert u.text == ingre["unit"]
    assert len(d.find_elements(By.XPATH, item % (i, ingre["item"]))) == 1
    assert len(d.find_elements(By.XPATH, quantity %
               (i, ingre["quantity"]))) == 1

  instr = list(map(
      lambda e: e.get_attribute("value"),
      d.find_elements(By.CSS_SELECTOR, ".recipe-instruction-input input")
  ))
  assert instr[:-1] == info["instr"]
