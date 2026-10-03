/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: rotateMediaAction
// ====================================================

export interface rotateMediaAction_rotateMedia {
  __typename: 'DevCmdResult'
  success: boolean
  message: string | null
}

export interface rotateMediaAction {
  /**
   * Rotate Media
   */
  rotateMedia: rotateMediaAction_rotateMedia
}

export interface rotateMediaActionVariables {
  mediaId: string
}
