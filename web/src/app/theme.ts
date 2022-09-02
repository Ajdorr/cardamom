export enum Theme {
  Focus,
  Surface,
  Background,
  Primary,
  PrimaryLight,
  PrimaryDark,
  Secondary,
  SecondaryLight,
  SecondaryDark,
  Tertiary,
  TertiaryLight,
  TertiaryDark,
  None,
}

export function themeToClass(theme: Theme): string {

  switch (theme) { 
    case Theme.Focus:
      return "theme-focus"
    case Theme.Surface:
      return "theme-surface"
    case Theme.Background:
      return "theme-background"
    default:
    case Theme.Primary:
      return "theme-primary"
    case Theme.PrimaryLight:
      return "theme-primary-light"
    case Theme.PrimaryDark:
      return "theme-primary-dark"
    case Theme.Secondary:
      return "theme-secondary"
    case Theme.SecondaryLight:
      return "theme-secondary-light"
    case Theme.SecondaryDark:
      return "theme-secondary-dark"
    case Theme.Tertiary:
      return "theme-tertiary"
    case Theme.TertiaryLight:
      return "theme-tertiary-light"
    case Theme.TertiaryDark:
      return "theme-tertiary-dark"
    case Theme.None:
      return ""
  }
}

export function joinWithClass(className: string, theme: Theme) {
  if (className) {
    return [className, themeToClass(theme)].join(" ")
  } else {
    return themeToClass(theme)
  }
}