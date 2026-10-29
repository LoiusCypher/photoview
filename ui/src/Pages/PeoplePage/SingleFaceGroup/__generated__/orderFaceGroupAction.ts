/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: orderFaceGroupAction
// ====================================================

export interface orderFaceGroupAction_orderFacesOfFaceGroup {
  __typename: 'DevCmdResult'
  success: boolean
  message: string | null
}

export interface orderFaceGroupAction {
  /**
   * Order unconfirmed faces in a FaceGroup
   */
  orderFacesOfFaceGroup: orderFaceGroupAction_orderFacesOfFaceGroup
}

export interface orderFaceGroupActionVariables {
  groupId: string
}
