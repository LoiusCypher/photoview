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
   * See if group contains several persons
   */
  splitFaceGroup: splitFaceGroupAction_splitFaceGroup
}

export interface splitFaceGroupActionVariables {
  groupId: string
}
