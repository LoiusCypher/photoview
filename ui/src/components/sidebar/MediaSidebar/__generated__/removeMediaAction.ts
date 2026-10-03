/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: removeMediaAction
// ====================================================

export interface removeMediaAction_removeMedia {
  __typename: 'DevCmdResult'
  success: boolean
  message: string | null
}

export interface removeMediaAction {
  /**
   * Remove Media
   */
  removeMedia: removeMediaAction_removeMedia
}

export interface removeMediaActionVariables {
  mediaId: string
}
