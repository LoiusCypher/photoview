/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: confirmFaceGroupAction
// ====================================================

export interface confirmFaceGroupAction_confirmAllFacesInFaceGroup {
  __typename: 'DevCmdResult'
  success: boolean
  message: string | null
}

export interface confirmFaceGroupAction {
  /**
   * Confirm all ImageFaces of given group
   */
  confirmAllFacesInFaceGroup: confirmFaceGroupAction_confirmAllFacesInFaceGroup
}

export interface confirmFaceGroupActionVariables {
  groupId: string
}
