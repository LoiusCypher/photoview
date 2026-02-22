/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: exportFaces
// ====================================================

export interface exportFaces_exportAllFaces {
  __typename: "ScannerResult";
  success: boolean;
  message: string | null;
}

export interface exportFaces {
  /**
   * Save al known portraits to dedicated folder
   */
  exportAllFaces: exportFaces_exportAllFaces;
}
