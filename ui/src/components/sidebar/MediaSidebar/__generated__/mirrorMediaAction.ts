/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: mirrorMediaAction
// ====================================================

export interface mirrorMediaAction_mirrorMedia {
  __typename: 'DevCmdResult'
  success: boolean
  message: string | null
}

export interface mirrorMediaAction {
  /**
   * Mirror Media
   */
  mirrorMedia: mirrorMediaAction_mirrorMedia
}

export interface mirrorMediaActionVariables {
  mediaId: string
}
