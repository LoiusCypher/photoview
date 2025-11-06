/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: watchModifiedQuery
// ====================================================

export interface watchModifiedQuery_siteInfo {
  __typename: "SiteInfo";
  /**
   * How often automatic scans should be initiated in seconds
   */
  watchModifiedTime: boolean;
}

export interface watchModifiedQuery {
  siteInfo: watchModifiedQuery_siteInfo;
}

