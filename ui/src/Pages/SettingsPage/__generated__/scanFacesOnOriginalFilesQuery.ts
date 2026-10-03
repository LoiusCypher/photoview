/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: scanFacesOnOriginalFilesQuery
// ====================================================

export interface scanFacesOnOriginalFilesQuery_siteInfo {
  __typename: 'SiteInfo'
  /**
   * If set original files are scanned for  faces insted of only thumbnails
   */
  scanFacesOnOriginalFiles: boolean
}

export interface scanFacesOnOriginalFilesQuery {
  siteInfo: scanFacesOnOriginalFilesQuery_siteInfo
}
