/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: splitFaceGroupAction
// ====================================================

export interface splitFaceGroupAction_splitFaceGroup {
  __typename: 'DevCmdResult'
  success: boolean
  message: string | null
}

export interface splitFaceGroupAction {
  /**
   * Try to isolate unconfirmed faces that could be split from current group
   */
  splitFaceGroup: splitFaceGroupAction_splitFaceGroup
}

export interface splitFaceGroupActionVariables {
  groupId: string
}
