/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: scanAlbumAction
// ====================================================

export interface scanAlbumAction_scanAlbum {
  __typename: 'ScannerResult'
  success: boolean
  message: string | null
}

export interface scanAlbumAction {
  /**
   * (Re-)Scan a single album for faces
   */
  scanAlbum: scanAlbumAction_scanAlbum
}

export interface scanAlbumActionVariables {
  albumId: string
}
