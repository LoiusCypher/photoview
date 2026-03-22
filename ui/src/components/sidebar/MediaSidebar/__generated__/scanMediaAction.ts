/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: scanMediaAction
// ====================================================

export interface scanMediaAction_scanMedia {
  __typename: 'ScannerResult'
  success: boolean
  message: string | null
}

export interface scanMediaAction {
  /**
   * (Re-)Scan a single media file for faces
   */
  scanMedia: scanMediaAction_scanMedia
}

export interface scanMediaActionVariables {
  mediaId: string
}
