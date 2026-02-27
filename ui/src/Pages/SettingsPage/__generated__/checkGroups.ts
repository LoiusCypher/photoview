/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: checkGroups
// ====================================================

export interface checkGroups_checkFaceGroup {
  __typename: "DevCmdResult";
  success: boolean;
  message: string | null;
}

export interface checkGroups {
  /**
   * Check FaceGroup consistency
   */
  checkFaceGroup: checkGroups_checkFaceGroup;
}

export interface checkGroupsVariables {
  faceGroupID: string;
}
