interface List {
  key: string;
  fiat: string;
  crypto: string;
  syntax: string;
}
interface EditForm {
  fiat: string;
  crypto: string;
  syntax: string;
}
export type { List, EditForm };
