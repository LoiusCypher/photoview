/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: reScanMediaAction
// ====================================================

export interface reScanMediaAction_reScanMedia {
  __typename: 'ScannerResult'
  success: boolean
  message: string | null
}

export interface reScanMediaAction {
  /**
   * ReScan a single media file for faces
   */
  reScanMedia: reScanMediaAction_reScanMedia
}

export interface reScanMediaActionVariables {
  mediaId: string
}
