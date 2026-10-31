/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: searchFaceGroupAction
// ====================================================

export interface searchFaceGroupAction_searchForFaceGroup {
  __typename: 'DevCmdResult'
  success: boolean
  message: string | null
}

export interface searchFaceGroupAction {
  /**
   * See if we could find faces for current group
   */
  searchForFaceGroup: searchFaceGroupAction_searchForFaceGroup
}

export interface searchFaceGroupActionVariables {
  groupId: string
}
