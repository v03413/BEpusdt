/* ArcoDesign命名空间 */
declare namespace ArcoDesign {
  interface ArcoSubmit {
    values: Record<string, any>;
    errors: Record<string, any> | undefined;
  }
}
