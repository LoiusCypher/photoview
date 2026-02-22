/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: scanMediaAction
// ====================================================

export interface scanMediaAction_scanMediaAction {
  __typename: "ScannerResult";
  success: boolean;
  message: string | null;
}

export interface scanMediaAction {
  /**
   * Scan a single media file for new media
   */
  scanMediaAction: scanMediaAction_scanMediaAction;
}

export interface scanMediaActionVariables {
  mediaId: string;
}
