/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: exportFaces
// ====================================================

export interface exportFaces_exportAllFaces {
  __typename: 'DevCmdResult'
  success: boolean
  message: string | null
}

export interface exportFaces {
  /**
   * Save all known portraits to dedicated folder
   */
  exportAllFaces: exportFaces_exportAllFaces
}
