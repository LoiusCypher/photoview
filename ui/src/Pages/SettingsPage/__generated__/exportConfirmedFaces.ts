/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: exportConfirmedFaces
// ====================================================

export interface exportConfirmedFaces_exportFaces {
  __typename: 'DevCmdResult'
  success: boolean
  message: string | null
}

export interface exportConfirmedFaces {
  /**
   * Save known portraits to dedicated folder
   */
  exportFaces: exportConfirmedFaces_exportFaces
}

export interface exportConfirmedFacesVariables {
  onlyConfirmed: boolean
}
