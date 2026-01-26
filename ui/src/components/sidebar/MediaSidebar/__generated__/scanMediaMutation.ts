/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: scanMediaMutation
// ====================================================

export interface scanMediaMutation_scanMediaAction {
  __typename: "ScannerResult";
  success: boolean;
  message: string | null;
}

export interface scanMediaMutation {
  /**
   * Scan a single media file for new media
   */
  scanMediaAction: scanMediaMutation_scanMediaAction;
}

export interface scanMediaMutationVariables {
  mediaId: string;
}
