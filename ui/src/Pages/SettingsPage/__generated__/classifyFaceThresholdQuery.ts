/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: classifyFaceThresholdQuery
// ====================================================

export interface classifyFaceThresholdQuery_siteInfo {
  __typename: 'SiteInfo'
  /**
   * Threshold for descriptor distance faces do not belong to same person
   */
  classifyFaceThreshold: number
}

export interface classifyFaceThresholdQuery {
  siteInfo: classifyFaceThresholdQuery_siteInfo
}
