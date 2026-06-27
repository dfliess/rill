/* eslint-disable */
import { getLocale, experimentalStaticLocale } from "../runtime.js"

/** @typedef {import('../runtime.js').LocalizedString} LocalizedString */
/** @typedef {{}} Alert_Context_Menu_AriaInputs */
/** @typedef {{ name: NonNullable<unknown> }} Alert_Created_ByInputs */
/** @typedef {{}} Alert_Created_Through_CodeInputs */
/** @typedef {{}} Alert_CriteriaInputs */
/** @typedef {{}} Alert_DashboardInputs */
/** @typedef {{}} Alert_DeleteInputs */
/** @typedef {{}} Alert_EditInputs */
/** @typedef {{}} Alert_Email_NotificationsInputs */
/** @typedef {{ count: NonNullable<unknown> }} Alert_Filters_LabelInputs */
/** @typedef {{}} Alert_Form_BackInputs */
/** @typedef {{}} Alert_Form_CancelInputs */
/** @typedef {{}} Alert_Form_CreateInputs */
/** @typedef {{}} Alert_Form_Create_TitleInputs */
/** @typedef {{}} Alert_Form_CreatedInputs */
/** @typedef {{}} Alert_Form_Criteria_DescriptionInputs */
/** @typedef {{}} Alert_Form_Criteria_Group_AriaInputs */
/** @typedef {{}} Alert_Form_Criteria_Measure_AriaInputs */
/** @typedef {{}} Alert_Form_Criteria_Measure_PlaceholderInputs */
/** @typedef {{}} Alert_Form_Criteria_Operator_AriaInputs */
/** @typedef {{}} Alert_Form_Criteria_Operator_PlaceholderInputs */
/** @typedef {{}} Alert_Form_Criteria_Preview_TitleInputs */
/** @typedef {{}} Alert_Form_Criteria_TitleInputs */
/** @typedef {{}} Alert_Form_Criteria_Type_AriaInputs */
/** @typedef {{}} Alert_Form_Criteria_Type_PlaceholderInputs */
/** @typedef {{}} Alert_Form_Criteria_Value_TitleInputs */
/** @typedef {{}} Alert_Form_Data_FiltersInputs */
/** @typedef {{}} Alert_Form_Data_MeasureInputs */
/** @typedef {{}} Alert_Form_Data_Measure_PlaceholderInputs */
/** @typedef {{}} Alert_Form_Data_Measures_DescInputs */
/** @typedef {{}} Alert_Form_Data_NoneInputs */
/** @typedef {{}} Alert_Form_Data_PreviewInputs */
/** @typedef {{}} Alert_Form_Data_Preview_DescInputs */
/** @typedef {{}} Alert_Form_Data_Split_ByInputs */
/** @typedef {{}} Alert_Form_Data_Split_PlaceholderInputs */
/** @typedef {{}} Alert_Form_Data_TitleInputs */
/** @typedef {{}} Alert_Form_Edit_TitleInputs */
/** @typedef {{}} Alert_Form_EditedInputs */
/** @typedef {{}} Alert_Form_Email_DescInputs */
/** @typedef {{}} Alert_Form_Email_PlaceholderInputs */
/** @typedef {{}} Alert_Form_Email_TitleInputs */
/** @typedef {{}} Alert_Form_Go_To_AlertsInputs */
/** @typedef {{}} Alert_Form_Name_PlaceholderInputs */
/** @typedef {{}} Alert_Form_Name_TitleInputs */
/** @typedef {{}} Alert_Form_NextInputs */
/** @typedef {{}} Alert_Form_No_CriteriaInputs */
/** @typedef {{}} Alert_Form_No_DataInputs */
/** @typedef {{}} Alert_Form_Preview_CellInputs */
/** @typedef {{}} Alert_Form_Preview_Table_AriaInputs */
/** @typedef {{}} Alert_Form_Select_CriteriaInputs */
/** @typedef {{}} Alert_Form_Slack_Channels_DescInputs */
/** @typedef {{ docsUrl: NonNullable<unknown> }} Alert_Form_Slack_Not_ConfiguredInputs */
/** @typedef {{}} Alert_Form_Slack_PlaceholderInputs */
/** @typedef {{}} Alert_Form_Slack_TitleInputs */
/** @typedef {{}} Alert_Form_Slack_Users_DescInputs */
/** @typedef {{}} Alert_Form_Snooze_DescInputs */
/** @typedef {{}} Alert_Form_Snooze_TitleInputs */
/** @typedef {{}} Alert_Form_Tab_CriteriaInputs */
/** @typedef {{}} Alert_Form_Tab_DataInputs */
/** @typedef {{}} Alert_Form_Tab_DeliveryInputs */
/** @typedef {{}} Alert_Form_TriggerInputs */
/** @typedef {{}} Alert_Form_Trigger_Data_RefreshInputs */
/** @typedef {{}} Alert_Form_Trigger_Set_ScheduleInputs */
/** @typedef {{}} Alert_Form_UpdateInputs */
/** @typedef {{ time: NonNullable<unknown> }} Alert_Last_CheckedInputs */
/** @typedef {{}} Alert_Name_LabelInputs */
/** @typedef {{}} Alert_No_Filters_BodyInputs */
/** @typedef {{}} Alert_No_Filters_HeadingInputs */
/** @typedef {{}} Alert_No_Filters_HintInputs */
/** @typedef {{}} Alert_NoneInputs */
/** @typedef {{}} Alert_Not_Checked_YetInputs */
/** @typedef {{}} Alert_ScheduleInputs */
/** @typedef {{}} Alert_Slack_NotificationsInputs */
/** @typedef {{}} Alert_SnoozeInputs */
/** @typedef {{}} Alert_Split_By_DimensionInputs */
/** @typedef {{}} Alert_Status_CheckedInputs */
/** @typedef {{}} Alert_Status_CheckingInputs */
/** @typedef {{}} Alert_Status_FailedInputs */
/** @typedef {{}} Alert_Status_Not_TriggeredInputs */
/** @typedef {{}} Alert_Status_RunningInputs */
/** @typedef {{}} Alert_Status_TriggeredInputs */
/** @typedef {{}} Alert_Status_UnknownInputs */
/** @typedef {{}} Alert_Unsubscribe_FailedInputs */
/** @typedef {{}} Alert_UnsubscribedInputs */
/** @typedef {{}} Alert_UnsubscribingInputs */
/** @typedef {{}} Alert_Whenever_Data_RefreshesInputs */
/** @typedef {{ alertsLink: NonNullable<unknown>, codeLink: NonNullable<unknown> }} Alerts_Empty_ActionInputs */
/** @typedef {{}} Alerts_Empty_MessageInputs */
/** @typedef {{}} Alerts_Link_TextInputs */
/** @typedef {{}} Alerts_Via_CodeInputs */
/** @typedef {{}} Avatar_Contact_SupportInputs */
/** @typedef {{}} Avatar_Copied_UrlInputs */
/** @typedef {{}} Avatar_Copy_UrlInputs */
/** @typedef {{}} Avatar_Copy_Url_For_ViewInputs */
/** @typedef {{}} Avatar_Create_Public_UrlInputs */
/** @typedef {{}} Avatar_DocumentationInputs */
/** @typedef {{}} Avatar_Join_DiscordInputs */
/** @typedef {{}} Avatar_LogoutInputs */
/** @typedef {{}} Avatar_ShareInputs */
/** @typedef {{}} Avatar_Share_DashboardInputs */
/** @typedef {{}} Avatar_Share_DescriptionInputs */
/** @typedef {{}} Avatar_View_AsInputs */
/** @typedef {{}} Bignumber_Copy_ValueInputs */
/** @typedef {{}} Bignumber_Shift_ClickInputs */
/** @typedef {{}} Bookmark_Absolute_Time_RangeInputs */
/** @typedef {{}} Bookmark_Absolute_Time_TooltipInputs */
/** @typedef {{}} Bookmark_CategoryInputs */
/** @typedef {{}} Bookmark_Category_TooltipInputs */
/** @typedef {{}} Bookmark_CreatedInputs */
/** @typedef {{}} Bookmark_Created_By_AdminInputs */
/** @typedef {{}} Bookmark_Current_ViewInputs */
/** @typedef {{}} Bookmark_Current_View_As_HomeInputs */
/** @typedef {{}} Bookmark_Default_LabelInputs */
/** @typedef {{}} Bookmark_Delete_BookmarkInputs */
/** @typedef {{}} Bookmark_Delete_Home_BookmarkInputs */
/** @typedef {{ name: NonNullable<unknown> }} Bookmark_DeletedInputs */
/** @typedef {{}} Bookmark_DescriptionInputs */
/** @typedef {{}} Bookmark_EditInputs */
/** @typedef {{}} Bookmark_FiltersInputs */
/** @typedef {{}} Bookmark_Filters_InheritedInputs */
/** @typedef {{}} Bookmark_Filters_Only_TooltipInputs */
/** @typedef {{}} Bookmark_Go_To_HomeInputs */
/** @typedef {{}} Bookmark_Home_CreatedInputs */
/** @typedef {{}} Bookmark_Home_DescriptionInputs */
/** @typedef {{}} Bookmark_LabelInputs */
/** @typedef {{}} Bookmark_Managed_BookmarksInputs */
/** @typedef {{}} Bookmark_No_BookmarksInputs */
/** @typedef {{}} Bookmark_No_SharedInputs */
/** @typedef {{}} Bookmark_Return_To_HomeInputs */
/** @typedef {{}} Bookmark_SaveInputs */
/** @typedef {{}} Bookmark_Save_Filters_OnlyInputs */
/** @typedef {{}} Bookmark_UpdatedInputs */
/** @typedef {{}} Bookmark_Your_BookmarksInputs */
/** @typedef {{}} Calendar_ApplyInputs */
/** @typedef {{ capLabel: NonNullable<unknown> }} Calendar_Range_Exceeds_LimitInputs */
/** @typedef {{}} Canvas_Add_WidgetInputs */
/** @typedef {{}} Canvas_Ai_Generating_ChartInputs */
/** @typedef {{}} Canvas_Ai_HintInputs */
/** @typedef {{}} Canvas_Ai_Is_EditingInputs */
/** @typedef {{}} Canvas_Ai_Manual_HintInputs */
/** @typedef {{}} Canvas_Ai_Write_ManuallyInputs */
/** @typedef {{}} Canvas_Align_BottomInputs */
/** @typedef {{}} Canvas_Align_CenterInputs */
/** @typedef {{}} Canvas_Align_LeftInputs */
/** @typedef {{}} Canvas_Align_MiddleInputs */
/** @typedef {{}} Canvas_Align_RightInputs */
/** @typedef {{}} Canvas_Align_TopInputs */
/** @typedef {{}} Canvas_Alignment_LabelInputs */
/** @typedef {{}} Canvas_Apply_Measure_Formatting_LabelInputs */
/** @typedef {{}} Canvas_Back_To_PromptInputs */
/** @typedef {{}} Canvas_BluesInputs */
/** @typedef {{}} Canvas_Breakdown_By_LabelInputs */
/** @typedef {{}} Canvas_CancelInputs */
/** @typedef {{}} Canvas_ChartInputs */
/** @typedef {{}} Canvas_Chart_TypeInputs */
/** @typedef {{}} Canvas_CividisInputs */
/** @typedef {{}} Canvas_Clear_FiltersInputs */
/** @typedef {{}} Canvas_Color_LabelInputs */
/** @typedef {{}} Canvas_Color_MappingInputs */
/** @typedef {{}} Canvas_Column_Dimensions_LabelInputs */
/** @typedef {{}} Canvas_Columns_LabelInputs */
/** @typedef {{}} Canvas_Comparison_Values_LabelInputs */
/** @typedef {{ label: NonNullable<unknown> }} Canvas_ConfigurationInputs */
/** @typedef {{}} Canvas_ConfigurationsInputs */
/** @typedef {{}} Canvas_Custom_ChartInputs */
/** @typedef {{}} Canvas_Data_Labels_LabelInputs */
/** @typedef {{}} Canvas_DeleteInputs */
/** @typedef {{}} Canvas_Delete_WidgetInputs */
/** @typedef {{}} Canvas_Delete_Widget_DescriptionInputs */
/** @typedef {{}} Canvas_Describe_Chart_ChangesInputs */
/** @typedef {{}} Canvas_Describe_Chart_PromptInputs */
/** @typedef {{}} Canvas_Description_LabelInputs */
/** @typedef {{}} Canvas_Description_PlaceholderInputs */
/** @typedef {{}} Canvas_Dimension_LabelInputs */
/** @typedef {{}} Canvas_Dimensions_LabelInputs */
/** @typedef {{}} Canvas_Display_NameInputs */
/** @typedef {{}} Canvas_Diverging_ThemeInputs */
/** @typedef {{}} Canvas_EditInputs */
/** @typedef {{}} Canvas_Edit_With_AiInputs */
/** @typedef {{}} Canvas_End_ColorInputs */
/** @typedef {{}} Canvas_Enter_A_NumberInputs */
/** @typedef {{}} Canvas_Evenly_DistributeInputs */
/** @typedef {{}} Canvas_Filter_BarInputs */
/** @typedef {{}} Canvas_FiltersInputs */
/** @typedef {{}} Canvas_FormatInputs */
/** @typedef {{}} Canvas_GenerateInputs */
/** @typedef {{}} Canvas_GradientInputs */
/** @typedef {{}} Canvas_GreensInputs */
/** @typedef {{}} Canvas_GreysInputs */
/** @typedef {{}} Canvas_Hide_Labels_BelowInputs */
/** @typedef {{}} Canvas_Hide_Total_Column_LabelInputs */
/** @typedef {{}} Canvas_Hide_Total_Row_LabelInputs */
/** @typedef {{}} Canvas_ImageInputs */
/** @typedef {{}} Canvas_InfernoInputs */
/** @typedef {{}} Canvas_Inner_Radius_LabelInputs */
/** @typedef {{}} Canvas_Insert_WidgetInputs */
/** @typedef {{ row: NonNullable<unknown>, col: NonNullable<unknown> }} Canvas_Insert_Widget_AtInputs */
/** @typedef {{ col: NonNullable<unknown> }} Canvas_Insert_Widget_At_ColInputs */
/** @typedef {{ row: NonNullable<unknown> }} Canvas_Insert_Widget_At_RowInputs */
/** @typedef {{}} Canvas_KpiInputs */
/** @typedef {{}} Canvas_Label_AngleInputs */
/** @typedef {{}} Canvas_LeaderboardInputs */
/** @typedef {{}} Canvas_Left_Y_Axis_LabelInputs */
/** @typedef {{}} Canvas_Legend_BottomInputs */
/** @typedef {{}} Canvas_Legend_LeftInputs */
/** @typedef {{}} Canvas_Legend_NoneInputs */
/** @typedef {{}} Canvas_Legend_OrientationInputs */
/** @typedef {{}} Canvas_Legend_RightInputs */
/** @typedef {{}} Canvas_Legend_TopInputs */
/** @typedef {{}} Canvas_LimitInputs */
/** @typedef {{}} Canvas_Local_FiltersInputs */
/** @typedef {{}} Canvas_Local_Time_RangeInputs */
/** @typedef {{}} Canvas_MagmaInputs */
/** @typedef {{}} Canvas_Mark_TypeInputs */
/** @typedef {{}} Canvas_Markdown_LabelInputs */
/** @typedef {{}} Canvas_MaxInputs */
/** @typedef {{}} Canvas_Max_WidthInputs */
/** @typedef {{}} Canvas_Measure_LabelInputs */
/** @typedef {{}} Canvas_Measures_LabelInputs */
/** @typedef {{}} Canvas_Metrics_Sql_LabelInputs */
/** @typedef {{}} Canvas_Metrics_View_DescriptionInputs */
/** @typedef {{}} Canvas_Metrics_View_LabelInputs */
/** @typedef {{}} Canvas_MinInputs */
/** @typedef {{}} Canvas_Mode_LabelInputs */
/** @typedef {{}} Canvas_Name_OptionInputs */
/** @typedef {{}} Canvas_No_Color_Values_FoundInputs */
/** @typedef {{}} Canvas_No_Components_AddedInputs */
/** @typedef {{}} Canvas_No_Filters_SelectedInputs */
/** @typedef {{ id: NonNullable<unknown> }} Canvas_No_Valid_ComponentInputs */
/** @typedef {{}} Canvas_No_Valid_Metrics_ViewInputs */
/** @typedef {{}} Canvas_Not_FoundInputs */
/** @typedef {{}} Canvas_Number_Of_Rows_LabelInputs */
/** @typedef {{}} Canvas_OrangesInputs */
/** @typedef {{}} Canvas_Order_OptionInputs */
/** @typedef {{}} Canvas_Percent_Of_LabelInputs */
/** @typedef {{}} Canvas_PlasmaInputs */
/** @typedef {{}} Canvas_Previous_OptionInputs */
/** @typedef {{}} Canvas_PurplesInputs */
/** @typedef {{}} Canvas_RedsInputs */
/** @typedef {{ idx: NonNullable<unknown> }} Canvas_Remove_Query_AriaInputs */
/** @typedef {{}} Canvas_Reset_To_DefaultInputs */
/** @typedef {{ row: NonNullable<unknown>, col: NonNullable<unknown> }} Canvas_Resize_AriaInputs */
/** @typedef {{}} Canvas_Right_Y_Axis_LabelInputs */
/** @typedef {{}} Canvas_Row_Dimensions_LabelInputs */
/** @typedef {{}} Canvas_Save_As_DefaultInputs */
/** @typedef {{}} Canvas_Saved_Default_FiltersInputs */
/** @typedef {{}} Canvas_Saving_Default_FiltersInputs */
/** @typedef {{}} Canvas_SchemeInputs */
/** @typedef {{}} Canvas_See_LessInputs */
/** @typedef {{}} Canvas_See_More_ValueInputs */
/** @typedef {{ count: NonNullable<unknown> }} Canvas_See_More_ValuesInputs */
/** @typedef {{}} Canvas_Select_Metrics_ViewInputs */
/** @typedef {{}} Canvas_Sequential_ThemeInputs */
/** @typedef {{}} Canvas_Show_Axis_TitleInputs */
/** @typedef {{}} Canvas_Show_Description_As_Tooltip_LabelInputs */
/** @typedef {{}} Canvas_Show_Null_ValuesInputs */
/** @typedef {{}} Canvas_Show_Other_Bucket_LabelInputs */
/** @typedef {{}} Canvas_Show_Totals_ValueInputs */
/** @typedef {{}} Canvas_Size_LabelInputs */
/** @typedef {{}} Canvas_SortInputs */
/** @typedef {{}} Canvas_Sparkline_BelowInputs */
/** @typedef {{}} Canvas_Sparkline_LabelInputs */
/** @typedef {{}} Canvas_Sparkline_RightInputs */
/** @typedef {{}} Canvas_SpectralInputs */
/** @typedef {{}} Canvas_Stage_LabelInputs */
/** @typedef {{}} Canvas_Start_ColorInputs */
/** @typedef {{ mark: NonNullable<unknown> }} Canvas_Switch_Mark_Type_AriaInputs */
/** @typedef {{}} Canvas_TableInputs */
/** @typedef {{}} Canvas_Table_TypeInputs */
/** @typedef {{}} Canvas_Teal_BluesInputs */
/** @typedef {{}} Canvas_TealsInputs */
/** @typedef {{}} Canvas_Text_MarkdownInputs */
/** @typedef {{}} Canvas_Time_Range_Display_LabelInputs */
/** @typedef {{}} Canvas_Time_RangesInputs */
/** @typedef {{}} Canvas_Time_ZonesInputs */
/** @typedef {{}} Canvas_Title_LabelInputs */
/** @typedef {{}} Canvas_Title_PlaceholderInputs */
/** @typedef {{}} Canvas_Tooltip_LabelInputs */
/** @typedef {{}} Canvas_Top_OptionInputs */
/** @typedef {{}} Canvas_TurboInputs */
/** @typedef {{}} Canvas_Unknown_ErrorInputs */
/** @typedef {{}} Canvas_Url_LabelInputs */
/** @typedef {{}} Canvas_Value_OptionInputs */
/** @typedef {{}} Canvas_Vega_Lite_Spec_LabelInputs */
/** @typedef {{}} Canvas_Viewing_Default_StateInputs */
/** @typedef {{}} Canvas_ViridisInputs */
/** @typedef {{}} Canvas_Which_Metrics_ViewInputs */
/** @typedef {{}} Canvas_Width_OptionInputs */
/** @typedef {{}} Canvas_X_Axis_LabelInputs */
/** @typedef {{}} Canvas_Y_Axis_LabelInputs */
/** @typedef {{}} Canvas_Zero_Based_OriginInputs */
/** @typedef {{}} Chart_AdaptiveInputs */
/** @typedef {{}} Chart_Adaptive_TooltipInputs */
/** @typedef {{ chartType: NonNullable<unknown> }} Chart_Add_Comparison_To_UseInputs */
/** @typedef {{}} Chart_BarInputs */
/** @typedef {{}} Chart_Copy_To_ClipboardInputs */
/** @typedef {{}} Chart_LineInputs */
/** @typedef {{}} Chart_Stacked_AreaInputs */
/** @typedef {{}} Chart_Stacked_BarInputs */
/** @typedef {{}} Chart_Undo_ZoomInputs */
/** @typedef {{}} Chart_Undo_Zoom_LabelInputs */
/** @typedef {{}} Chart_VsInputs */
/** @typedef {{}} Chart_ZoomInputs */
/** @typedef {{}} Chart_Zoom_LabelInputs */
/** @typedef {{}} Chat_Ai_DisclaimerInputs */
/** @typedef {{}} Chat_Cancel_StreamingInputs */
/** @typedef {{}} Chat_CloseInputs */
/** @typedef {{}} Chat_Connect_ClientInputs */
/** @typedef {{}} Chat_Conversation_HistoryInputs */
/** @typedef {{}} Chat_Downvote_AriaInputs */
/** @typedef {{}} Chat_Downvote_TooltipInputs */
/** @typedef {{}} Chat_Duration_Less_Than_SecondInputs */
/** @typedef {{ count: NonNullable<unknown> }} Chat_Duration_MinutesInputs */
/** @typedef {{}} Chat_Duration_One_MinuteInputs */
/** @typedef {{}} Chat_Duration_One_SecondInputs */
/** @typedef {{ count: NonNullable<unknown> }} Chat_Duration_SecondsInputs */
/** @typedef {{}} Chat_Empty_LabelInputs */
/** @typedef {{}} Chat_Failed_To_GenerateInputs */
/** @typedef {{}} Chat_Feedback_AnalyzingInputs */
/** @typedef {{}} Chat_Feedback_CommentsInputs */
/** @typedef {{}} Chat_Feedback_PlaceholderInputs */
/** @typedef {{}} Chat_Feedback_Select_AllInputs */
/** @typedef {{}} Chat_Feedback_SkipInputs */
/** @typedef {{}} Chat_Feedback_SubmitInputs */
/** @typedef {{}} Chat_Feedback_TitleInputs */
/** @typedef {{ days: NonNullable<unknown> }} Chat_Group_Days_AgoInputs */
/** @typedef {{}} Chat_Group_OlderInputs */
/** @typedef {{}} Chat_Group_TodayInputs */
/** @typedef {{}} Chat_Group_YesterdayInputs */
/** @typedef {{}} Chat_Happy_To_ExploreInputs */
/** @typedef {{}} Chat_How_Can_I_HelpInputs */
/** @typedef {{}} Chat_New_ConversationInputs */
/** @typedef {{}} Chat_No_ConversationsInputs */
/** @typedef {{}} Chat_Placeholder_AnalystInputs */
/** @typedef {{}} Chat_Send_MessageInputs */
/** @typedef {{}} Chat_Share_ConversationInputs */
/** @typedef {{}} Chat_Share_CopiedInputs */
/** @typedef {{}} Chat_Share_Create_LinkInputs */
/** @typedef {{}} Chat_Share_CreatingInputs */
/** @typedef {{}} Chat_Share_DescriptionInputs */
/** @typedef {{}} Chat_Share_Start_FirstInputs */
/** @typedef {{}} Chat_Show_DetailsInputs */
/** @typedef {{}} Chat_ThinkingInputs */
/** @typedef {{ duration: NonNullable<unknown> }} Chat_Thought_ForInputs */
/** @typedef {{}} Chat_Unable_To_LoadInputs */
/** @typedef {{}} Chat_Upvote_AriaInputs */
/** @typedef {{}} Chat_Upvote_TooltipInputs */
/** @typedef {{}} Common_ApplyInputs */
/** @typedef {{}} Common_CancelInputs */
/** @typedef {{}} Common_Close_SearchInputs */
/** @typedef {{}} Common_ContinueInputs */
/** @typedef {{}} Common_Must_Be_NumberInputs */
/** @typedef {{}} Common_PreviewInputs */
/** @typedef {{}} Common_RequiredInputs */
/** @typedef {{}} Common_SearchInputs */
/** @typedef {{}} Common_Search_ListInputs */
/** @typedef {{}} Common_UndoInputs */
/** @typedef {{ name: NonNullable<unknown> }} Dashboard_Add_All_Dims_RowsInputs */
/** @typedef {{}} Dashboard_Add_All_To_ColumnsInputs */
/** @typedef {{ name: NonNullable<unknown> }} Dashboard_Add_All_To_Columns_TagInputs */
/** @typedef {{}} Dashboard_Add_All_To_RowsInputs */
/** @typedef {{}} Dashboard_Add_ColumnInputs */
/** @typedef {{}} Dashboard_Add_FilterInputs */
/** @typedef {{}} Dashboard_Add_Filter_ButtonInputs */
/** @typedef {{}} Dashboard_Add_Filter_Button_AriaInputs */
/** @typedef {{}} Dashboard_Add_Mock_UserInputs */
/** @typedef {{}} Dashboard_Add_RowInputs */
/** @typedef {{}} Dashboard_Add_To_ColumnsInputs */
/** @typedef {{}} Dashboard_Add_To_RowsInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Added_Items_FilterInputs */
/** @typedef {{}} Dashboard_AllInputs */
/** @typedef {{}} Dashboard_All_MeasuresInputs */
/** @typedef {{}} Dashboard_Always_Show_AsInputs */
/** @typedef {{}} Dashboard_Anchor_Period_EndInputs */
/** @typedef {{}} Dashboard_As_OfInputs */
/** @typedef {{}} Dashboard_As_Of_RefInputs */
/** @typedef {{}} Dashboard_Auto_ArrangeInputs */
/** @typedef {{ name: NonNullable<unknown> }} Dashboard_Auto_Arrange_TagInputs */
/** @typedef {{}} Dashboard_CancelInputs */
/** @typedef {{}} Dashboard_Change_Over_ComparisonInputs */
/** @typedef {{}} Dashboard_ClearInputs */
/** @typedef {{}} Dashboard_Clear_FiltersInputs */
/** @typedef {{}} Dashboard_Clear_RecentsInputs */
/** @typedef {{}} Dashboard_Clear_ViewInputs */
/** @typedef {{}} Dashboard_Click_To_Edit_ValuesInputs */
/** @typedef {{}} Dashboard_Cmd_Click_ReplaceInputs */
/** @typedef {{}} Dashboard_Collapse_AllInputs */
/** @typedef {{}} Dashboard_ColumnsInputs */
/** @typedef {{}} Dashboard_Comparison_Column_AriaInputs */
/** @typedef {{}} Dashboard_Comparison_Select_TooltipInputs */
/** @typedef {{}} Dashboard_Complete_DataInputs */
/** @typedef {{}} Dashboard_Complete_Data_DescriptionInputs */
/** @typedef {{}} Dashboard_Connect_Sparse_DataInputs */
/** @typedef {{}} Dashboard_ContainsInputs */
/** @typedef {{}} Dashboard_Copied_Error_ClipboardInputs */
/** @typedef {{}} Dashboard_Copy_ErrorInputs */
/** @typedef {{}} Dashboard_Copy_Error_ClipboardInputs */
/** @typedef {{}} Dashboard_Copy_Error_To_Clipboard_AriaInputs */
/** @typedef {{}} Dashboard_Current_TimeInputs */
/** @typedef {{}} Dashboard_Current_Time_DescriptionInputs */
/** @typedef {{}} Dashboard_Currently_FlatInputs */
/** @typedef {{}} Dashboard_Currently_PivotInputs */
/** @typedef {{}} Dashboard_CustomInputs */
/** @typedef {{}} Dashboard_Deselect_AllInputs */
/** @typedef {{}} Dashboard_Deselect_All_SelectionsInputs */
/** @typedef {{}} Dashboard_Dimension_Display_AriaInputs */
/** @typedef {{}} Dashboard_Dimension_Search_Results_AriaInputs */
/** @typedef {{}} Dashboard_Dimension_Table_AriaInputs */
/** @typedef {{}} Dashboard_DimensionsInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Dimensions_CountInputs */
/** @typedef {{}} Dashboard_Dimensions_LabelInputs */
/** @typedef {{}} Dashboard_Download_PngInputs */
/** @typedef {{}} Dashboard_Drag_DimensionsInputs */
/** @typedef {{}} Dashboard_Drag_Dimensions_Or_MeasuresInputs */
/** @typedef {{ zone: NonNullable<unknown> }} Dashboard_Drag_List_ZoneInputs */
/** @typedef {{}} Dashboard_Dynamic_Y_AxisInputs */
/** @typedef {{}} Dashboard_EndInputs */
/** @typedef {{}} Dashboard_Error_OccurredInputs */
/** @typedef {{}} Dashboard_Error_Occurred_HoverInputs */
/** @typedef {{}} Dashboard_Error_TagInputs */
/** @typedef {{}} Dashboard_Errored_Contact_AdminInputs */
/** @typedef {{ link: NonNullable<unknown> }} Dashboard_Errored_Need_HelpInputs */
/** @typedef {{}} Dashboard_Errored_TitleInputs */
/** @typedef {{}} Dashboard_Errored_View_ProjectInputs */
/** @typedef {{}} Dashboard_Errored_View_StatusInputs */
/** @typedef {{}} Dashboard_Errored_View_Status_ButtonInputs */
/** @typedef {{}} Dashboard_ExcludeInputs */
/** @typedef {{}} Dashboard_ExploreInputs */
/** @typedef {{}} Dashboard_Export_ChartInputs */
/** @typedef {{}} Dashboard_Export_Dimension_Table_DataInputs */
/** @typedef {{}} Dashboard_Export_Model_DataInputs */
/** @typedef {{}} Dashboard_Export_Pivot_DataInputs */
/** @typedef {{}} Dashboard_Export_Table_DataInputs */
/** @typedef {{}} Dashboard_Filter_By_ValueInputs */
/** @typedef {{}} Dashboard_Filter_Dimension_ValueInputs */
/** @typedef {{}} Dashboard_Filter_Required_Set_ValueInputs */
/** @typedef {{}} Dashboard_FlatInputs */
/** @typedef {{ time: NonNullable<unknown> }} Dashboard_GeneratedInputs */
/** @typedef {{}} Dashboard_GeneratingInputs */
/** @typedef {{}} Dashboard_GrainInputs */
/** @typedef {{}} Dashboard_Hide_PanelsInputs */
/** @typedef {{}} Dashboard_In_ListInputs */
/** @typedef {{}} Dashboard_Include_Exclude_ToggleInputs */
/** @typedef {{}} Dashboard_Invalid_Time_RangeInputs */
/** @typedef {{ time: NonNullable<unknown> }} Dashboard_Last_Refreshed_AgoInputs */
/** @typedef {{}} Dashboard_Latest_DataInputs */
/** @typedef {{}} Dashboard_Latest_Data_DescriptionInputs */
/** @typedef {{}} Dashboard_Leaderboards_AriaInputs */
/** @typedef {{}} Dashboard_List_CreateInputs */
/** @typedef {{ link: NonNullable<unknown> }} Dashboard_List_Create_To_StartInputs */
/** @typedef {{}} Dashboard_List_EmptyInputs */
/** @typedef {{}} Dashboard_List_See_AllInputs */
/** @typedef {{ name: NonNullable<unknown> }} Dashboard_Measure_Chart_AriaInputs */
/** @typedef {{}} Dashboard_MeasuresInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Measures_CountInputs */
/** @typedef {{}} Dashboard_Measures_LabelInputs */
/** @typedef {{}} Dashboard_Menu_All_DimensionsInputs */
/** @typedef {{}} Dashboard_No_Additional_DetailsInputs */
/** @typedef {{}} Dashboard_No_Available_FieldsInputs */
/** @typedef {{}} Dashboard_No_Comparison_DimensionInputs */
/** @typedef {{}} Dashboard_No_Filters_SelectedInputs */
/** @typedef {{}} Dashboard_No_Matching_TagsInputs */
/** @typedef {{}} Dashboard_No_Mock_UsersInputs */
/** @typedef {{}} Dashboard_No_Options_FoundInputs */
/** @typedef {{}} Dashboard_No_Search_ResultsInputs */
/** @typedef {{}} Dashboard_No_Timezones_ConfiguredInputs */
/** @typedef {{}} Dashboard_No_Valid_GrainsInputs */
/** @typedef {{}} Dashboard_OfInputs */
/** @typedef {{}} Dashboard_Open_Dimension_Details_AriaInputs */
/** @typedef {{}} Dashboard_OtherInputs */
/** @typedef {{}} Dashboard_OthersInputs */
/** @typedef {{}} Dashboard_Output_ExcludesInputs */
/** @typedef {{}} Dashboard_Output_IncludesInputs */
/** @typedef {{}} Dashboard_Percent_Of_TotalInputs */
/** @typedef {{}} Dashboard_Percentage_ChangeInputs */
/** @typedef {{}} Dashboard_PivotInputs */
/** @typedef {{}} Dashboard_Pivot_Add_MeasureInputs */
/** @typedef {{}} Dashboard_Pivot_Building_TableInputs */
/** @typedef {{}} Dashboard_Pivot_Give_DataInputs */
/** @typedef {{}} Dashboard_Pivot_Keep_It_UpInputs */
/** @typedef {{}} Dashboard_Pivot_Learn_MoreInputs */
/** @typedef {{}} Dashboard_Pivot_Need_Help_DiscordInputs */
/** @typedef {{}} Dashboard_Pivot_No_DataInputs */
/** @typedef {{}} Dashboard_Pivot_Table_LonelyInputs */
/** @typedef {{}} Dashboard_Readonly_Filter_Chips_AriaInputs */
/** @typedef {{}} Dashboard_RecentInputs */
/** @typedef {{}} Dashboard_ReferenceInputs */
/** @typedef {{ label: NonNullable<unknown> }} Dashboard_Remove_LabelInputs */
/** @typedef {{ count: NonNullable<unknown> }} Dashboard_Removed_Items_FilterInputs */
/** @typedef {{}} Dashboard_ReplaceInputs */
/** @typedef {{ name: NonNullable<unknown> }} Dashboard_Replace_Auto_ArrangeInputs */
/** @typedef {{ name: NonNullable<unknown> }} Dashboard_Replace_Columns_TagInputs */
/** @typedef {{}} Dashboard_Replace_Columns_Tag_ItemsInputs */
/** @typedef {{}} Dashboard_Replace_Pivot_DescriptionInputs */
/** @typedef {{}} Dashboard_Replace_Pivot_TitleInputs */
/** @typedef {{}} Dashboard_Replace_Rows_Cols_TagInputs */
/** @typedef {{ name: NonNullable<unknown> }} Dashboard_Replace_Rows_TagInputs */
/** @typedef {{}} Dashboard_Replace_Rows_Tag_DimsInputs */
/** @typedef {{}} Dashboard_Required_MeasureInputs */
/** @typedef {{}} Dashboard_Row_LimitInputs */
/** @typedef {{}} Dashboard_Row_Limit_TooltipInputs */
/** @typedef {{}} Dashboard_RowsInputs */
/** @typedef {{}} Dashboard_Search_DimensionInputs */
/** @typedef {{}} Dashboard_Search_DimensionsInputs */
/** @typedef {{}} Dashboard_Search_ResultsInputs */
/** @typedef {{}} Dashboard_See_MoreInputs */
/** @typedef {{}} Dashboard_Select_Aggregation_Grain_AriaInputs */
/** @typedef {{}} Dashboard_Select_AllInputs */
/** @typedef {{}} Dashboard_Select_Comparison_DimensionInputs */
/** @typedef {{}} Dashboard_Select_Comparison_HintInputs */
/** @typedef {{}} Dashboard_Select_Ref_Time_GrainInputs */
/** @typedef {{}} Dashboard_Select_Time_AxisInputs */
/** @typedef {{}} Dashboard_Select_Time_Comparison_AriaInputs */
/** @typedef {{ label: NonNullable<unknown> }} Dashboard_Select_Time_DimensionInputs */
/** @typedef {{}} Dashboard_Select_Time_GrainInputs */
/** @typedef {{}} Dashboard_Select_Time_RangeInputs */
/** @typedef {{}} Dashboard_Select_Time_Range_AriaInputs */
/** @typedef {{}} Dashboard_SelectedInputs */
/** @typedef {{}} Dashboard_Show_PanelsInputs */
/** @typedef {{}} Dashboard_Sort_By_Absolute_Change_AriaInputs */
/** @typedef {{}} Dashboard_Sort_By_Percent_Change_AriaInputs */
/** @typedef {{}} Dashboard_Sort_By_Percent_Total_AriaInputs */
/** @typedef {{}} Dashboard_Sort_By_Value_AriaInputs */
/** @typedef {{}} Dashboard_StartInputs */
/** @typedef {{}} Dashboard_Start_PivotInputs */
/** @typedef {{}} Dashboard_Switch_FlatInputs */
/** @typedef {{}} Dashboard_Switch_PivotInputs */
/** @typedef {{}} Dashboard_Table_ModeInputs */
/** @typedef {{}} Dashboard_TagsInputs */
/** @typedef {{}} Dashboard_Tdd_Contact_DiscordInputs */
/** @typedef {{}} Dashboard_Tdd_ErrorInputs */
/** @typedef {{}} Dashboard_Tdd_No_ComparisonInputs */
/** @typedef {{}} Dashboard_Tdd_TimeInputs */
/** @typedef {{}} Dashboard_TimeInputs */
/** @typedef {{}} Dashboard_Time_AxisInputs */
/** @typedef {{}} Dashboard_Time_LabelInputs */
/** @typedef {{}} Dashboard_Time_ZoneInputs */
/** @typedef {{}} Dashboard_Timezone_SelectorInputs */
/** @typedef {{}} Dashboard_ToInputs */
/** @typedef {{}} Dashboard_Toggle_ExcludeInputs */
/** @typedef {{}} Dashboard_Toggle_IncludeInputs */
/** @typedef {{}} Dashboard_Toggle_Rows_Viewer_AriaInputs */
/** @typedef {{}} Dashboard_Toggle_Time_Comparison_AriaInputs */
/** @typedef {{}} Dashboard_Total_ColumnInputs */
/** @typedef {{}} Dashboard_Total_RowInputs */
/** @typedef {{}} Dashboard_View_AsInputs */
/** @typedef {{}} Dashboard_Viewing_AsInputs */
/** @typedef {{}} Dialog_Close_Without_Saving_Alert_DescInputs */
/** @typedef {{}} Dialog_Close_Without_Saving_CancelInputs */
/** @typedef {{}} Dialog_Close_Without_Saving_ConfirmInputs */
/** @typedef {{}} Dialog_Close_Without_Saving_TitleInputs */
/** @typedef {{}} Error_Access_Denied_BodyInputs */
/** @typedef {{}} Error_Access_Denied_HeaderInputs */
/** @typedef {{}} Error_Auth_BodyInputs */
/** @typedef {{}} Error_Auth_HeaderInputs */
/** @typedef {{}} Error_Back_To_HomeInputs */
/** @typedef {{}} Error_Conversation_Not_Found_BodyInputs */
/** @typedef {{}} Error_Conversation_Not_Found_HeaderInputs */
/** @typedef {{}} Error_Deploying_ProjectInputs */
/** @typedef {{}} Error_Deployment_ErrorInputs */
/** @typedef {{}} Error_Deployment_Not_Found_BodyInputs */
/** @typedef {{}} Error_Deployment_Not_Found_HeaderInputs */
/** @typedef {{}} Error_Fetching_DeploymentInputs */
/** @typedef {{}} Error_Generic_BodyInputs */
/** @typedef {{}} Error_Generic_HeaderInputs */
/** @typedef {{}} Error_Hide_DetailsInputs */
/** @typedef {{}} Error_Link_Expired_BodyInputs */
/** @typedef {{}} Error_Link_Expired_HeaderInputs */
/** @typedef {{}} Error_Network_BodyInputs */
/** @typedef {{}} Error_Network_HeaderInputs */
/** @typedef {{}} Error_Org_Not_Found_BodyInputs */
/** @typedef {{}} Error_Org_Not_Found_HeaderInputs */
/** @typedef {{}} Error_Page_Not_Found_BodyInputs */
/** @typedef {{}} Error_Page_Not_Found_HeaderInputs */
/** @typedef {{}} Error_Project_Not_Found_BodyInputs */
/** @typedef {{}} Error_Project_Not_Found_HeaderInputs */
/** @typedef {{}} Error_Resource_Not_Found_BodyInputs */
/** @typedef {{}} Error_Resource_Not_Found_HeaderInputs */
/** @typedef {{}} Error_Retry_NowInputs */
/** @typedef {{}} Error_Show_DetailsInputs */
/** @typedef {{}} Explore_All_DimensionsInputs */
/** @typedef {{}} Explore_All_MeasuresInputs */
/** @typedef {{}} Explore_By_Grain_PrefixInputs */
/** @typedef {{}} Explore_Choose_DimensionsInputs */
/** @typedef {{}} Explore_Choose_MeasuresInputs */
/** @typedef {{}} Explore_Clear_FilterInputs */
/** @typedef {{ tag: NonNullable<unknown> }} Explore_Clear_Filter_TagInputs */
/** @typedef {{}} Explore_Clear_Search_To_Reorder_DimensionsInputs */
/** @typedef {{}} Explore_Clear_Search_To_Reorder_MeasuresInputs */
/** @typedef {{}} Explore_Clear_Tag_FilterInputs */
/** @typedef {{}} Explore_Clear_Tag_Filter_To_Reorder_DimensionsInputs */
/** @typedef {{}} Explore_Clear_Tag_Filter_To_Reorder_MeasuresInputs */
/** @typedef {{ count: NonNullable<unknown>, total: NonNullable<unknown> }} Explore_Dimensions_CountInputs */
/** @typedef {{ tag: NonNullable<unknown> }} Explore_Filter_By_TagInputs */
/** @typedef {{}} Explore_Filter_LabelInputs */
/** @typedef {{}} Explore_Go_To_DashboardInputs */
/** @typedef {{}} Explore_Go_To_ExploreInputs */
/** @typedef {{ name: NonNullable<unknown> }} Explore_Go_To_NamedInputs */
/** @typedef {{}} Explore_Hidden_DimensionsInputs */
/** @typedef {{}} Explore_Hidden_MeasuresInputs */
/** @typedef {{}} Explore_Hide_AllInputs */
/** @typedef {{ tag: NonNullable<unknown> }} Explore_Hide_All_In_Named_TagInputs */
/** @typedef {{}} Explore_Hide_All_In_TagInputs */
/** @typedef {{ name: NonNullable<unknown> }} Explore_Hide_ItemInputs */
/** @typedef {{ count: NonNullable<unknown>, total: NonNullable<unknown> }} Explore_Measures_CountInputs */
/** @typedef {{}} Explore_Multi_SelectInputs */
/** @typedef {{}} Explore_Must_Show_One_DimensionInputs */
/** @typedef {{}} Explore_Must_Show_One_MeasureInputs */
/** @typedef {{ count: NonNullable<unknown> }} Explore_N_MeasuresInputs */
/** @typedef {{}} Explore_No_Dimensions_From_TagInputs */
/** @typedef {{}} Explore_No_Dimensions_Or_TagsInputs */
/** @typedef {{}} Explore_No_Dimensions_ShownInputs */
/** @typedef {{}} Explore_No_Hidden_DimensionsInputs */
/** @typedef {{}} Explore_No_Hidden_MeasuresInputs */
/** @typedef {{}} Explore_No_Matching_Dimensions_ShownInputs */
/** @typedef {{}} Explore_No_Matching_Hidden_DimensionsInputs */
/** @typedef {{}} Explore_No_Matching_Hidden_MeasuresInputs */
/** @typedef {{}} Explore_No_Matching_Leaderboard_MeasuresInputs */
/** @typedef {{}} Explore_No_Matching_Measures_ShownInputs */
/** @typedef {{}} Explore_No_Matching_TagsInputs */
/** @typedef {{}} Explore_No_Measures_From_TagInputs */
/** @typedef {{}} Explore_No_Measures_Or_TagsInputs */
/** @typedef {{}} Explore_No_Measures_ShownInputs */
/** @typedef {{ tag: NonNullable<unknown> }} Explore_Only_Show_TagInputs */
/** @typedef {{}} Explore_Only_Show_This_TagInputs */
/** @typedef {{}} Explore_Search_DimensionsInputs */
/** @typedef {{}} Explore_Search_Dimensions_Or_TagsInputs */
/** @typedef {{}} Explore_Search_ListInputs */
/** @typedef {{}} Explore_Search_MeasuresInputs */
/** @typedef {{}} Explore_Search_Measures_Or_TagsInputs */
/** @typedef {{}} Explore_Show_AllInputs */
/** @typedef {{ tag: NonNullable<unknown> }} Explore_Show_All_In_Named_TagInputs */
/** @typedef {{}} Explore_Show_All_In_TagInputs */
/** @typedef {{}} Explore_Show_Context_For_All_MeasuresInputs */
/** @typedef {{ name: NonNullable<unknown> }} Explore_Show_ItemInputs */
/** @typedef {{}} Explore_ShowingInputs */
/** @typedef {{}} Explore_Shown_DimensionsInputs */
/** @typedef {{}} Explore_Shown_MeasuresInputs */
/** @typedef {{ visible: NonNullable<unknown>, total: NonNullable<unknown> }} Explore_Tag_Shown_CountInputs */
/** @typedef {{}} Explore_TagsInputs */
/** @typedef {{}} Explore_Unable_To_OpenInputs */
/** @typedef {{}} Explore_Unknown_DimensionInputs */
/** @typedef {{}} Explore_Unknown_MeasureInputs */
/** @typedef {{ label: NonNullable<unknown> }} Field_List_Add_FieldsInputs */
/** @typedef {{ label: NonNullable<unknown> }} Field_List_AriaInputs */
/** @typedef {{}} Field_List_DimensionsInputs */
/** @typedef {{}} Field_List_MeasuresInputs */
/** @typedef {{}} Field_List_TimeInputs */
/** @typedef {{}} Filter_Advanced_BetaInputs */
/** @typedef {{}} Filter_Advanced_WarningInputs */
/** @typedef {{}} Filter_DimensionsInputs */
/** @typedef {{}} Filter_Enter_Search_TermInputs */
/** @typedef {{}} Filter_Make_OptionalInputs */
/** @typedef {{}} Filter_Make_RequiredInputs */
/** @typedef {{ dimension: NonNullable<unknown> }} Filter_Measure_For_DimensionInputs */
/** @typedef {{ comparison: NonNullable<unknown> }} Filter_Measure_From_ComparisonInputs */
/** @typedef {{}} Filter_Measure_Op_BetweenInputs */
/** @typedef {{}} Filter_Measure_Op_Does_Not_EqualInputs */
/** @typedef {{}} Filter_Measure_Op_EqualsInputs */
/** @typedef {{}} Filter_Measure_Op_Greater_ThanInputs */
/** @typedef {{}} Filter_Measure_Op_Greater_Than_Or_EqualsInputs */
/** @typedef {{}} Filter_Measure_Op_Less_ThanInputs */
/** @typedef {{}} Filter_Measure_Op_Less_Than_Or_EqualsInputs */
/** @typedef {{}} Filter_Measure_Op_Not_BetweenInputs */
/** @typedef {{}} Filter_Measure_Type_Change_FromInputs */
/** @typedef {{}} Filter_Measure_Type_Percent_Change_FromInputs */
/** @typedef {{}} Filter_Measure_Type_Percent_Of_TotalInputs */
/** @typedef {{}} Filter_Measure_Type_ValueInputs */
/** @typedef {{}} Filter_MeasuresInputs */
/** @typedef {{}} Filter_Mode_ContainsInputs */
/** @typedef {{}} Filter_Mode_Contains_DescriptionInputs */
/** @typedef {{}} Filter_Mode_In_ListInputs */
/** @typedef {{}} Filter_Mode_In_List_DescriptionInputs */
/** @typedef {{}} Filter_Mode_SelectInputs */
/** @typedef {{}} Filter_Mode_Select_DescriptionInputs */
/** @typedef {{}} Filter_Paste_List_HintInputs */
/** @typedef {{}} Filter_PinInputs */
/** @typedef {{}} Filter_Pin_TooltipInputs */
/** @typedef {{}} Filter_Required_TooltipInputs */
/** @typedef {{}} Filter_UnpinInputs */
/** @typedef {{}} Footer_Report_IssueInputs */
/** @typedef {{}} Footer_Rill_DeveloperInputs */
/** @typedef {{}} Footer_Shortcut_ClickInputs */
/** @typedef {{}} Footer_Unknown_VersionInputs */
/** @typedef {{}} Footer_VersionInputs */
/** @typedef {{}} Footer_View_DocumentationInputs */
/** @typedef {{}} Form_OptionalInputs */
/** @typedef {{}} Groups_Changes_SavedInputs */
/** @typedef {{}} Groups_Create_A_GroupInputs */
/** @typedef {{}} Groups_Create_GroupInputs */
/** @typedef {{}} Groups_CreatedInputs */
/** @typedef {{}} Groups_DeleteInputs */
/** @typedef {{}} Groups_Delete_Confirm_DescInputs */
/** @typedef {{}} Groups_Delete_Confirm_TitleInputs */
/** @typedef {{}} Groups_DeletedInputs */
/** @typedef {{}} Groups_EditInputs */
/** @typedef {{}} Groups_Edit_GroupInputs */
/** @typedef {{}} Groups_Error_Adding_RoleInputs */
/** @typedef {{}} Groups_Error_DeletingInputs */
/** @typedef {{}} Groups_Error_LoadingInputs */
/** @typedef {{}} Groups_Error_Revoking_RoleInputs */
/** @typedef {{}} Groups_Error_Updating_RoleInputs */
/** @typedef {{}} Groups_RemovedInputs */
/** @typedef {{}} Groups_RenamedInputs */
/** @typedef {{}} Groups_Role_AddedInputs */
/** @typedef {{}} Groups_Role_RevokedInputs */
/** @typedef {{}} Groups_Role_UpdatedInputs */
/** @typedef {{}} Groups_Table_EmptyInputs */
/** @typedef {{}} Groups_Table_Header_GroupInputs */
/** @typedef {{ count: NonNullable<unknown> }} Groups_Total_CountInputs */
/** @typedef {{}} Groups_Yes_DeleteInputs */
/** @typedef {{}} Home_Dashboards_HeadingInputs */
/** @typedef {{}} Home_Subtitle_No_ChatInputs */
/** @typedef {{}} Home_Subtitle_With_ChatInputs */
/** @typedef {{ projectName: NonNullable<unknown> }} Home_Welcome_ToInputs */
/** @typedef {{}} Interval_DayInputs */
/** @typedef {{}} Interval_HourInputs */
/** @typedef {{}} Interval_NoneInputs */
/** @typedef {{}} Interval_WeekInputs */
/** @typedef {{}} Kpi_No_ChangeInputs */
/** @typedef {{}} Kpi_No_DataInputs */
/** @typedef {{}} Kpi_Not_AvailableInputs */
/** @typedef {{ comparison: NonNullable<unknown> }} Kpi_Vs_ComparisonInputs */
/** @typedef {{}} Language_EnInputs */
/** @typedef {{}} Language_EsInputs */
/** @typedef {{}} Language_Switcher_LabelInputs */
/** @typedef {{}} Layout_Inspector_Panel_AriaInputs */
/** @typedef {{}} Leaderboard_CompareInputs */
/** @typedef {{}} Leaderboard_Copy_ValueInputs */
/** @typedef {{}} Leaderboard_Expand_TableInputs */
/** @typedef {{}} Leaderboard_Expand_TooltipInputs */
/** @typedef {{}} Leaderboard_No_Available_ValuesInputs */
/** @typedef {{}} Leaderboard_Remove_ComparisonInputs */
/** @typedef {{}} Leaderboard_Shift_ClickInputs */
/** @typedef {{ name: NonNullable<unknown> }} Leaderboard_Toggle_BreakdownInputs */
/** @typedef {{}} Mcp_Add_To_ConfigInputs */
/** @typedef {{}} Mcp_Add_UrlInputs */
/** @typedef {{}} Mcp_ConfigurationInputs */
/** @typedef {{}} Mcp_Create_TokenInputs */
/** @typedef {{ privateLabel: NonNullable<unknown>, tokenLabel: NonNullable<unknown> }} Mcp_Create_Token_DescInputs */
/** @typedef {{}} Mcp_Create_Token_TitleInputs */
/** @typedef {{}} Mcp_Dialog_DescriptionInputs */
/** @typedef {{}} Mcp_Dialog_TitleInputs */
/** @typedef {{}} Mcp_IssuingInputs */
/** @typedef {{}} Mcp_Learn_MoreInputs */
/** @typedef {{}} Mcp_Manual_TabInputs */
/** @typedef {{}} Mcp_Oauth_AutoInputs */
/** @typedef {{}} Mcp_Oauth_TabInputs */
/** @typedef {{}} Mcp_Personal_Access_TokenInputs */
/** @typedef {{}} Mcp_PrivateInputs */
/** @typedef {{}} Mcp_RecommendedInputs */
/** @typedef {{}} Mcp_Token_CreatedInputs */
/** @typedef {{}} Mcp_Token_FailedInputs */
/** @typedef {{}} Measure_Filter_ApplyInputs */
/** @typedef {{}} Measure_Filter_By_DimensionInputs */
/** @typedef {{}} Measure_Filter_Enter_NumberInputs */
/** @typedef {{}} Measure_Filter_Higher_ValueInputs */
/** @typedef {{}} Measure_Filter_Lower_ValueInputs */
/** @typedef {{}} Measure_Filter_Select_DimensionInputs */
/** @typedef {{}} Measure_Filter_ThresholdInputs */
/** @typedef {{}} Measures_Choose_TooltipInputs */
/** @typedef {{}} Measures_LabelInputs */
/** @typedef {{}} Nav_Close_SidebarInputs */
/** @typedef {{}} Nav_Data_ExplorerInputs */
/** @typedef {{}} Nav_Show_SidebarInputs */
/** @typedef {{}} Nav_Tab_AiInputs */
/** @typedef {{}} Nav_Tab_AlertsInputs */
/** @typedef {{}} Nav_Tab_DashboardsInputs */
/** @typedef {{}} Nav_Tab_HomeInputs */
/** @typedef {{}} Nav_Tab_QueryInputs */
/** @typedef {{}} Nav_Tab_ReportsInputs */
/** @typedef {{}} Nav_Tab_SettingsInputs */
/** @typedef {{}} Nav_Tab_StatusInputs */
/** @typedef {{}} Org_Check_Out_ProjectsInputs */
/** @typedef {{}} Org_New_ProjectInputs */
/** @typedef {{}} Org_Search_Add_Remove_UsersInputs */
/** @typedef {{}} Org_Tab_ProjectsInputs */
/** @typedef {{}} Org_Tab_SettingsInputs */
/** @typedef {{}} Org_Tab_UsersInputs */
/** @typedef {{}} Pivot_Collapse_RowInputs */
/** @typedef {{}} Pivot_Dim_OneInputs */
/** @typedef {{}} Pivot_Dim_OtherInputs */
/** @typedef {{}} Pivot_Drop_Arrange_AriaInputs */
/** @typedef {{}} Pivot_Drop_Replace_AriaInputs */
/** @typedef {{ dimCount: NonNullable<unknown>, dimLabel: NonNullable<unknown>, measureCount: NonNullable<unknown>, measureLabel: NonNullable<unknown> }} Pivot_Drop_Replace_TextInputs */
/** @typedef {{}} Pivot_Drop_Split_HintInputs */
/** @typedef {{ dimCount: NonNullable<unknown>, dimLabel: NonNullable<unknown>, measureCount: NonNullable<unknown>, measureLabel: NonNullable<unknown> }} Pivot_Drop_Split_TextInputs */
/** @typedef {{}} Pivot_Expand_RowInputs */
/** @typedef {{}} Pivot_Measure_OneInputs */
/** @typedef {{}} Pivot_Measure_OtherInputs */
/** @typedef {{}} Pivot_TagsInputs */
/** @typedef {{ grain: NonNullable<unknown> }} Pivot_Time_Dimension_HeaderInputs */
/** @typedef {{}} Pivot_Time_PrefixInputs */
/** @typedef {{}} Project_Dashboards_TitleInputs */
/** @typedef {{}} Project_DeleteInputs */
/** @typedef {{}} Project_EditInputs */
/** @typedef {{}} Project_RenameInputs */
/** @typedef {{}} Project_Role_AdminInputs */
/** @typedef {{}} Project_Role_ViewerInputs */
/** @typedef {{}} Project_Search_Or_InviteInputs */
/** @typedef {{}} Project_Search_UsersInputs */
/** @typedef {{}} Project_ShareInputs */
/** @typedef {{ project: NonNullable<unknown> }} Project_Share_HeadingInputs */
/** @typedef {{}} Project_Share_TooltipInputs */
/** @typedef {{}} Report_Context_Menu_AriaInputs */
/** @typedef {{ name: NonNullable<unknown> }} Report_Created_ByInputs */
/** @typedef {{}} Report_Created_Through_CodeInputs */
/** @typedef {{}} Report_DashboardInputs */
/** @typedef {{}} Report_DeleteInputs */
/** @typedef {{}} Report_EditInputs */
/** @typedef {{}} Report_Email_RecipientsInputs */
/** @typedef {{}} Report_Email_ValidationInputs */
/** @typedef {{}} Report_Form_CancelInputs */
/** @typedef {{}} Report_Form_ChannelsInputs */
/** @typedef {{}} Report_Form_Clear_FiltersInputs */
/** @typedef {{}} Report_Form_ColumnsInputs */
/** @typedef {{}} Report_Form_CreateInputs */
/** @typedef {{}} Report_Form_Create_ButtonInputs */
/** @typedef {{}} Report_Form_Created_NotificationInputs */
/** @typedef {{}} Report_Form_DayInputs */
/** @typedef {{}} Report_Form_Day_FirstInputs */
/** @typedef {{}} Report_Form_Day_FridayInputs */
/** @typedef {{}} Report_Form_Day_MondayInputs */
/** @typedef {{}} Report_Form_Day_SaturdayInputs */
/** @typedef {{}} Report_Form_Day_SundayInputs */
/** @typedef {{}} Report_Form_Day_ThursdayInputs */
/** @typedef {{}} Report_Form_Day_TuesdayInputs */
/** @typedef {{}} Report_Form_Day_WednesdayInputs */
/** @typedef {{}} Report_Form_DocsInputs */
/** @typedef {{}} Report_Form_Edited_NotificationInputs */
/** @typedef {{}} Report_Form_Email_HintInputs */
/** @typedef {{}} Report_Form_Email_PlaceholderInputs */
/** @typedef {{}} Report_Form_Email_RecipientsInputs */
/** @typedef {{}} Report_Form_Email_RecurringInputs */
/** @typedef {{}} Report_Form_FiltersInputs */
/** @typedef {{}} Report_Form_Filters_AriaInputs */
/** @typedef {{}} Report_Form_FormatInputs */
/** @typedef {{}} Report_Form_Format_CsvInputs */
/** @typedef {{}} Report_Form_Format_ParquetInputs */
/** @typedef {{}} Report_Form_Format_XlsxInputs */
/** @typedef {{}} Report_Form_Freq_DailyInputs */
/** @typedef {{}} Report_Form_Freq_MonthlyInputs */
/** @typedef {{}} Report_Form_Freq_WeekdaysInputs */
/** @typedef {{}} Report_Form_Freq_WeeklyInputs */
/** @typedef {{}} Report_Form_FrequencyInputs */
/** @typedef {{}} Report_Form_Go_To_ReportsInputs */
/** @typedef {{}} Report_Form_Include_MetadataInputs */
/** @typedef {{}} Report_Form_Invalid_EmailInputs */
/** @typedef {{}} Report_Form_Local_SuffixInputs */
/** @typedef {{}} Report_Form_Metadata_TooltipInputs */
/** @typedef {{}} Report_Form_No_Rows_SelectedInputs */
/** @typedef {{}} Report_Form_Recipients_Must_Be_ProjectInputs */
/** @typedef {{}} Report_Form_RequiredInputs */
/** @typedef {{}} Report_Form_Row_LimitInputs */
/** @typedef {{}} Report_Form_Row_Limit_PlaceholderInputs */
/** @typedef {{}} Report_Form_RowsInputs */
/** @typedef {{}} Report_Form_Run_AsInputs */
/** @typedef {{}} Report_Form_Run_As_CreatorInputs */
/** @typedef {{}} Report_Form_Run_As_Creator_DescInputs */
/** @typedef {{}} Report_Form_Run_As_RecipientInputs */
/** @typedef {{}} Report_Form_Run_As_Recipient_DescInputs */
/** @typedef {{}} Report_Form_SaveInputs */
/** @typedef {{}} Report_Form_Save_ButtonInputs */
/** @typedef {{}} Report_Form_ScheduleInputs */
/** @typedef {{}} Report_Form_Slack_Channels_HintInputs */
/** @typedef {{ link: NonNullable<unknown> }} Report_Form_Slack_Not_ConfiguredInputs */
/** @typedef {{}} Report_Form_Slack_TitleInputs */
/** @typedef {{}} Report_Form_Slack_UsersInputs */
/** @typedef {{}} Report_Form_Slack_Users_HintInputs */
/** @typedef {{}} Report_Form_TimeInputs */
/** @typedef {{}} Report_Form_TimezoneInputs */
/** @typedef {{}} Report_Form_Title_LabelInputs */
/** @typedef {{}} Report_Form_Title_PlaceholderInputs */
/** @typedef {{}} Report_Name_LabelInputs */
/** @typedef {{}} Report_Next_RunInputs */
/** @typedef {{}} Report_No_Row_LimitInputs */
/** @typedef {{}} Report_RepeatsInputs */
/** @typedef {{}} Report_Slack_RecipientsInputs */
/** @typedef {{}} Report_Status_FailedInputs */
/** @typedef {{}} Report_Status_SentInputs */
/** @typedef {{}} Report_Triggered_AdhocInputs */
/** @typedef {{}} Report_Unsubscribe_FailedInputs */
/** @typedef {{}} Report_UnsubscribedInputs */
/** @typedef {{}} Report_UnsubscribingInputs */
/** @typedef {{ reportsLink: NonNullable<unknown> }} Reports_Empty_ActionInputs */
/** @typedef {{}} Reports_Empty_MessageInputs */
/** @typedef {{}} Reports_Link_TextInputs */
/** @typedef {{}} Resource_Error_Contact_SupportInputs */
/** @typedef {{ kind: NonNullable<unknown> }} Resource_Error_LoadingInputs */
/** @typedef {{}} Resource_Search_PlaceholderInputs */
/** @typedef {{}} Resource_Type_CanvasInputs */
/** @typedef {{}} Resource_Type_ExploreInputs */
/** @typedef {{}} Role_AdminInputs */
/** @typedef {{}} Role_EditorInputs */
/** @typedef {{}} Role_GuestInputs */
/** @typedef {{}} Role_Guest_DescInputs */
/** @typedef {{}} Role_Org_Admin_DescInputs */
/** @typedef {{}} Role_Org_Admin_DetailInputs */
/** @typedef {{}} Role_Org_Editor_DescInputs */
/** @typedef {{}} Role_Org_Editor_DetailInputs */
/** @typedef {{}} Role_Org_Viewer_DescInputs */
/** @typedef {{}} Role_Org_Viewer_DetailInputs */
/** @typedef {{}} Role_Project_Admin_DescInputs */
/** @typedef {{}} Role_Project_Editor_DescInputs */
/** @typedef {{}} Role_Project_Viewer_DescInputs */
/** @typedef {{}} Role_ViewerInputs */
/** @typedef {{ link: NonNullable<unknown> }} Share_Limited_ViewInputs */
/** @typedef {{}} Share_Original_DashboardInputs */
/** @typedef {{}} Snooze_OffInputs */
/** @typedef {{}} Status_Action_DescribeInputs */
/** @typedef {{}} Status_Action_Full_RefreshInputs */
/** @typedef {{}} Status_Action_Incremental_RefreshInputs */
/** @typedef {{}} Status_Action_Refresh_Errored_PartitionsInputs */
/** @typedef {{}} Status_Action_View_LogsInputs */
/** @typedef {{}} Status_Action_View_PartitionsInputs */
/** @typedef {{}} Status_CancelInputs */
/** @typedef {{}} Status_Checking_ErrorsInputs */
/** @typedef {{}} Status_ClearInputs */
/** @typedef {{}} Status_Clone_DescriptionInputs */
/** @typedef {{}} Status_Column_Database_SizeInputs */
/** @typedef {{}} Status_Column_Last_RefreshInputs */
/** @typedef {{}} Status_Column_Model_NameInputs */
/** @typedef {{}} Status_Column_NameInputs */
/** @typedef {{}} Status_Column_Next_RefreshInputs */
/** @typedef {{}} Status_Column_Table_NameInputs */
/** @typedef {{}} Status_Column_TypeInputs */
/** @typedef {{}} Status_CompleteInputs */
/** @typedef {{}} Status_Compute_UnitInputs */
/** @typedef {{}} Status_Compute_UnitsInputs */
/** @typedef {{}} Status_Data_AccessibleInputs */
/** @typedef {{}} Status_Data_SizeInputs */
/** @typedef {{}} Status_Deploy_DeletedInputs */
/** @typedef {{}} Status_Deploy_DeletingInputs */
/** @typedef {{}} Status_Deploy_ErrorInputs */
/** @typedef {{}} Status_Deploy_Not_DeployedInputs */
/** @typedef {{}} Status_Deploy_PendingInputs */
/** @typedef {{}} Status_Deploy_ReadyInputs */
/** @typedef {{}} Status_Deploy_StoppedInputs */
/** @typedef {{}} Status_Deploy_StoppingInputs */
/** @typedef {{}} Status_Deploy_UpdatingInputs */
/** @typedef {{}} Status_DeploymentInputs */
/** @typedef {{}} Status_Download_ProjectInputs */
/** @typedef {{}} Status_Error_Loading_ResourcesInputs */
/** @typedef {{}} Status_Error_Loading_TablesInputs */
/** @typedef {{}} Status_ErrorsInputs */
/** @typedef {{}} Status_External_Tables_SectionInputs */
/** @typedef {{}} Status_Filter_ErrorInputs */
/** @typedef {{}} Status_Filter_OkInputs */
/** @typedef {{}} Status_Filter_WarnInputs */
/** @typedef {{}} Status_Full_Refresh_WarningInputs */
/** @typedef {{}} Status_Incremental_Refresh_DescriptionInputs */
/** @typedef {{}} Status_Label_Ai_ConnectorInputs */
/** @typedef {{}} Status_Label_BranchInputs */
/** @typedef {{}} Status_Label_Cluster_SizeInputs */
/** @typedef {{}} Status_Label_EnvironmentInputs */
/** @typedef {{}} Status_Label_Last_SyncedInputs */
/** @typedef {{}} Status_Label_Olap_EngineInputs */
/** @typedef {{}} Status_Label_RepoInputs */
/** @typedef {{}} Status_Label_RuntimeInputs */
/** @typedef {{}} Status_Label_StatusInputs */
/** @typedef {{}} Status_Learn_About_External_OlapInputs */
/** @typedef {{}} Status_Learn_MoreInputs */
/** @typedef {{}} Status_Load_More_TablesInputs */
/** @typedef {{}} Status_LoadingInputs */
/** @typedef {{}} Status_Loading_ModelsInputs */
/** @typedef {{}} Status_Loading_ResourcesInputs */
/** @typedef {{}} Status_Loading_TablesInputs */
/** @typedef {{}} Status_Logs_All_LevelsInputs */
/** @typedef {{}} Status_Logs_ConnectingInputs */
/** @typedef {{}} Status_Logs_Connection_FailedInputs */
/** @typedef {{}} Status_Logs_DisconnectedInputs */
/** @typedef {{}} Status_Logs_IdleInputs */
/** @typedef {{}} Status_Logs_Level_DebugInputs */
/** @typedef {{}} Status_Logs_Level_ErrorInputs */
/** @typedef {{}} Status_Logs_Level_InfoInputs */
/** @typedef {{}} Status_Logs_Level_WarnInputs */
/** @typedef {{}} Status_Logs_LiveInputs */
/** @typedef {{}} Status_Logs_No_MatchInputs */
/** @typedef {{}} Status_Logs_RetryInputs */
/** @typedef {{}} Status_Logs_WaitingInputs */
/** @typedef {{}} Status_Model_PartitionsInputs */
/** @typedef {{}} Status_Models_Created_In_DeveloperInputs */
/** @typedef {{}} Status_Models_SectionInputs */
/** @typedef {{}} Status_Nav_AnalyticsInputs */
/** @typedef {{}} Status_Nav_BranchesInputs */
/** @typedef {{}} Status_Nav_LogsInputs */
/** @typedef {{}} Status_Nav_OverviewInputs */
/** @typedef {{}} Status_Nav_ResourcesInputs */
/** @typedef {{}} Status_Nav_TablesInputs */
/** @typedef {{}} Status_No_ErrorsInputs */
/** @typedef {{}} Status_No_External_TablesInputs */
/** @typedef {{}} Status_No_External_Tables_Match_FiltersInputs */
/** @typedef {{}} Status_No_ModelsInputs */
/** @typedef {{}} Status_No_Models_Match_FiltersInputs */
/** @typedef {{}} Status_No_Parse_ErrorsInputs */
/** @typedef {{}} Status_No_Resource_DataInputs */
/** @typedef {{}} Status_No_ResourcesInputs */
/** @typedef {{}} Status_No_Resources_Match_FiltersInputs */
/** @typedef {{}} Status_No_TablesInputs */
/** @typedef {{}} Status_NoteInputs */
/** @typedef {{}} Status_Owned_By_OtherInputs */
/** @typedef {{}} Status_Owned_By_YouInputs */
/** @typedef {{}} Status_Page_TitleInputs */
/** @typedef {{}} Status_Parse_ErrorInputs */
/** @typedef {{}} Status_Parse_ErrorsInputs */
/** @typedef {{}} Status_Parse_Errors_TitleInputs */
/** @typedef {{}} Status_ReconcilingInputs */
/** @typedef {{}} Status_Refresh_AllInputs */
/** @typedef {{}} Status_Refresh_All_Confirm_BodyInputs */
/** @typedef {{}} Status_Refresh_All_Confirm_TipInputs */
/** @typedef {{}} Status_Refresh_All_Confirm_TitleInputs */
/** @typedef {{}} Status_Refresh_All_Sources_ModelsInputs */
/** @typedef {{}} Status_Refresh_Errored_Confirm_BodyInputs */
/** @typedef {{ modelName: NonNullable<unknown> }} Status_Refresh_Errored_Confirm_TitleInputs */
/** @typedef {{}} Status_RefreshingInputs */
/** @typedef {{}} Status_Resource_ReconcilingInputs */
/** @typedef {{}} Status_Rill_ManagedInputs */
/** @typedef {{}} Status_Table_PluralInputs */
/** @typedef {{}} Status_Table_SingularInputs */
/** @typedef {{}} Status_Unable_To_Check_ErrorsInputs */
/** @typedef {{}} Status_View_AllInputs */
/** @typedef {{}} Status_View_PluralInputs */
/** @typedef {{}} Status_View_SingularInputs */
/** @typedef {{}} Status_Yes_RefreshInputs */
/** @typedef {{}} Theme_DarkInputs */
/** @typedef {{}} Theme_LabelInputs */
/** @typedef {{}} Theme_LightInputs */
/** @typedef {{}} Theme_SystemInputs */
/** @typedef {{}} Time_1_Day_AgoInputs */
/** @typedef {{}} Time_1_Hour_AgoInputs */
/** @typedef {{}} Time_1_Minute_AgoInputs */
/** @typedef {{}} Time_1_Month_AgoInputs */
/** @typedef {{}} Time_1_Week_AgoInputs */
/** @typedef {{}} Time_1_Year_AgoInputs */
/** @typedef {{ duration: NonNullable<unknown> }} Time_AgoInputs */
/** @typedef {{}} Time_All_TimeInputs */
/** @typedef {{}} Time_ComparingInputs */
/** @typedef {{}} Time_Comparison_Previous_DayInputs */
/** @typedef {{}} Time_Comparison_Previous_PeriodInputs */
/** @typedef {{}} Time_CustomInputs */
/** @typedef {{}} Time_Custom_RangeInputs */
/** @typedef {{}} Time_Enter_Time_RangeInputs */
/** @typedef {{ duration: NonNullable<unknown> }} Time_From_NowInputs */
/** @typedef {{}} Time_Grain_ByInputs */
/** @typedef {{}} Time_Grain_CompleteInputs */
/** @typedef {{}} Time_Grain_DayInputs */
/** @typedef {{}} Time_Grain_DaysInputs */
/** @typedef {{}} Time_Grain_HourInputs */
/** @typedef {{}} Time_Grain_HoursInputs */
/** @typedef {{}} Time_Grain_MinuteInputs */
/** @typedef {{}} Time_Grain_MinutesInputs */
/** @typedef {{}} Time_Grain_MonthInputs */
/** @typedef {{}} Time_Grain_MonthsInputs */
/** @typedef {{}} Time_Grain_QuarterInputs */
/** @typedef {{}} Time_Grain_QuartersInputs */
/** @typedef {{}} Time_Grain_TimeInputs */
/** @typedef {{}} Time_Grain_WeekInputs */
/** @typedef {{}} Time_Grain_WeeksInputs */
/** @typedef {{}} Time_Grain_YearInputs */
/** @typedef {{}} Time_Grain_YearsInputs */
/** @typedef {{}} Time_Just_NowInputs */
/** @typedef {{ duration: NonNullable<unknown> }} Time_Last_DurationInputs */
/** @typedef {{}} Time_Month_To_DateInputs */
/** @typedef {{ count: NonNullable<unknown> }} Time_N_Days_AgoInputs */
/** @typedef {{ count: NonNullable<unknown> }} Time_N_Hours_AgoInputs */
/** @typedef {{ count: NonNullable<unknown> }} Time_N_Minutes_AgoInputs */
/** @typedef {{ count: NonNullable<unknown> }} Time_N_Months_AgoInputs */
/** @typedef {{ count: NonNullable<unknown> }} Time_N_Weeks_AgoInputs */
/** @typedef {{ count: NonNullable<unknown> }} Time_N_Years_AgoInputs */
/** @typedef {{}} Time_No_Comparison_DimensionInputs */
/** @typedef {{}} Time_No_Comparison_PeriodInputs */
/** @typedef {{}} Time_Previous_MonthInputs */
/** @typedef {{}} Time_Previous_QuarterInputs */
/** @typedef {{}} Time_Previous_WeekInputs */
/** @typedef {{}} Time_Previous_YearInputs */
/** @typedef {{}} Time_Quarter_To_DateInputs */
/** @typedef {{ grain: NonNullable<unknown> }} Time_Range_Grain_To_DateInputs */
/** @typedef {{ count: NonNullable<unknown>, grains: NonNullable<unknown> }} Time_Range_Last_N_GrainsInputs */
/** @typedef {{ grain: NonNullable<unknown> }} Time_Range_Next_GrainInputs */
/** @typedef {{ count: NonNullable<unknown>, grains: NonNullable<unknown> }} Time_Range_Next_N_GrainsInputs */
/** @typedef {{ grain: NonNullable<unknown> }} Time_Range_Previous_GrainInputs */
/** @typedef {{ grain: NonNullable<unknown> }} Time_Range_This_GrainInputs */
/** @typedef {{}} Time_Ref_CompleteInputs */
/** @typedef {{}} Time_Ref_Complete_DataInputs */
/** @typedef {{}} Time_Ref_CurrentInputs */
/** @typedef {{}} Time_Ref_EarliestInputs */
/** @typedef {{}} Time_Ref_LatestInputs */
/** @typedef {{}} Time_Ref_NowInputs */
/** @typedef {{ count: NonNullable<unknown> }} Time_Relative_Hours_ShortInputs */
/** @typedef {{ count: NonNullable<unknown> }} Time_Relative_Minutes_ShortInputs */
/** @typedef {{}} Time_Relative_NowInputs */
/** @typedef {{}} Time_TodayInputs */
/** @typedef {{}} Time_Unable_To_ParseInputs */
/** @typedef {{}} Time_VsInputs */
/** @typedef {{}} Time_Week_To_DateInputs */
/** @typedef {{}} Time_Year_To_DateInputs */
/** @typedef {{}} Time_YesterdayInputs */
/** @typedef {{}} Users_Access_Changed_EveryoneInputs */
/** @typedef {{}} Users_Access_Changed_Invite_OnlyInputs */
/** @typedef {{}} Users_Access_LevelInputs */
/** @typedef {{}} Users_Add_GuestInputs */
/** @typedef {{}} Users_Add_Guests_ButtonInputs */
/** @typedef {{}} Users_Add_Guests_DescriptionInputs */
/** @typedef {{}} Users_Add_Guests_TitleInputs */
/** @typedef {{}} Users_Add_UsersInputs */
/** @typedef {{ count: NonNullable<unknown> }} Users_Added_Groups_CountInputs */
/** @typedef {{ emails: NonNullable<unknown> }} Users_Already_MemberInputs */
/** @typedef {{ count: NonNullable<unknown> }} Users_And_MoreInputs */
/** @typedef {{}} Users_Billing_Change_Role_DescInputs */
/** @typedef {{}} Users_Billing_Change_Role_TitleInputs */
/** @typedef {{}} Users_Billing_Remove_DescInputs */
/** @typedef {{}} Users_Billing_Remove_TitleInputs */
/** @typedef {{}} Users_CancelInputs */
/** @typedef {{}} Users_Change_Billing_ContactInputs */
/** @typedef {{}} Users_ConvertInputs */
/** @typedef {{}} Users_Convert_To_MemberInputs */
/** @typedef {{ user: NonNullable<unknown>, role: NonNullable<unknown> }} Users_Convert_User_To_RoleInputs */
/** @typedef {{}} Users_Copy_UrlInputs */
/** @typedef {{}} Users_CreateInputs */
/** @typedef {{}} Users_Email_Or_Group_PlaceholderInputs */
/** @typedef {{}} Users_Email_PlaceholderInputs */
/** @typedef {{}} Users_ErrorInputs */
/** @typedef {{}} Users_Error_Loading_MembersInputs */
/** @typedef {{}} Users_Error_RemovingInputs */
/** @typedef {{}} Users_Error_Removing_UserInputs */
/** @typedef {{}} Users_Error_Updating_RoleInputs */
/** @typedef {{}} Users_Error_Upgrading_RoleInputs */
/** @typedef {{ organization: NonNullable<unknown> }} Users_Everyone_At_OrgInputs */
/** @typedef {{ groups: NonNullable<unknown> }} Users_Failed_Add_GroupsInputs */
/** @typedef {{ emails: NonNullable<unknown> }} Users_Failed_InviteInputs */
/** @typedef {{ emails: NonNullable<unknown> }} Users_Failed_Invite_UsersInputs */
/** @typedef {{}} Users_Failed_Load_ProjectsInputs */
/** @typedef {{}} Users_Failed_Load_UsersInputs */
/** @typedef {{}} Users_Filter_AdminsInputs */
/** @typedef {{}} Users_Filter_AllInputs */
/** @typedef {{}} Users_Filter_All_RolesInputs */
/** @typedef {{}} Users_Filter_All_UsersInputs */
/** @typedef {{}} Users_Filter_EditorsInputs */
/** @typedef {{}} Users_Filter_MembersInputs */
/** @typedef {{}} Users_Filter_Pending_InvitesInputs */
/** @typedef {{}} Users_Filter_ViewersInputs */
/** @typedef {{}} Users_Form_NameInputs */
/** @typedef {{}} Users_Form_UntitledInputs */
/** @typedef {{}} Users_Form_UsersInputs */
/** @typedef {{}} Users_General_AccessInputs */
/** @typedef {{ count: NonNullable<unknown> }} Users_Group_CountInputs */
/** @typedef {{ role: NonNullable<unknown> }} Users_Guest_UpgradedInputs */
/** @typedef {{ role: NonNullable<unknown> }} Users_Guest_Upgraded_ToInputs */
/** @typedef {{ count: NonNullable<unknown>, role: NonNullable<unknown> }} Users_Guests_Invited_SuccessInputs */
/** @typedef {{}} Users_InviteInputs */
/** @typedef {{}} Users_Invite_OnlyInputs */
/** @typedef {{ count: NonNullable<unknown> }} Users_Invited_CountInputs */
/** @typedef {{ count: NonNullable<unknown>, role: NonNullable<unknown> }} Users_Invited_SuccessInputs */
/** @typedef {{}} Users_Learn_More_SharingInputs */
/** @typedef {{}} Users_LoadingInputs */
/** @typedef {{}} Users_Loading_MoreInputs */
/** @typedef {{ count: NonNullable<unknown> }} Users_Member_CountInputs */
/** @typedef {{}} Users_No_GroupsInputs */
/** @typedef {{}} Users_No_ProjectsInputs */
/** @typedef {{}} Users_No_UsersInputs */
/** @typedef {{}} Users_Only_Admins_AccessInputs */
/** @typedef {{}} Users_Org_Members_AccessInputs */
/** @typedef {{}} Users_Page_TitleInputs */
/** @typedef {{}} Users_Pending_InvitationInputs */
/** @typedef {{}} Users_Project_AccessInputs */
/** @typedef {{ count: NonNullable<unknown> }} Users_Project_CountInputs */
/** @typedef {{ count: NonNullable<unknown> }} Users_Projects_CountInputs */
/** @typedef {{}} Users_RemoveInputs */
/** @typedef {{}} Users_Remove_Confirm_DescInputs */
/** @typedef {{}} Users_Remove_Confirm_TitleInputs */
/** @typedef {{}} Users_RemovedInputs */
/** @typedef {{}} Users_Removed_From_OrgInputs */
/** @typedef {{}} Users_RetryInputs */
/** @typedef {{}} Users_Role_UpdatedInputs */
/** @typedef {{}} Users_SaveInputs */
/** @typedef {{}} Users_SearchingInputs */
/** @typedef {{}} Users_Section_GroupsInputs */
/** @typedef {{}} Users_Section_GuestsInputs */
/** @typedef {{}} Users_Section_MembersInputs */
/** @typedef {{}} Users_Select_ProjectsInputs */
/** @typedef {{ count: NonNullable<unknown> }} Users_Tab_GroupsInputs */
/** @typedef {{ count: NonNullable<unknown> }} Users_Tab_GuestsInputs */
/** @typedef {{ count: NonNullable<unknown> }} Users_Tab_MembersInputs */
/** @typedef {{}} Users_Table_EmptyInputs */
/** @typedef {{}} Users_Table_Header_GroupsInputs */
/** @typedef {{}} Users_Table_Header_Org_RoleInputs */
/** @typedef {{}} Users_Table_Header_ProjectsInputs */
/** @typedef {{}} Users_Table_Header_UserInputs */
/** @typedef {{ role: NonNullable<unknown> }} Users_Upgrade_Confirm_DescInputs */
/** @typedef {{ role: NonNullable<unknown> }} Users_Upgrade_Confirm_TitleInputs */
/** @typedef {{}} Users_Url_CopiedInputs */
/** @typedef {{ count: NonNullable<unknown> }} Users_User_CountInputs */
/** @typedef {{}} Users_Yes_RemoveInputs */
/** @typedef {{}} Users_Yes_UpgradeInputs */
import * as __en from "./en.js"
import * as __es from "./es.js"
/**
* | output |
* | --- |
* | "Alert context menu" |
*
* @param {Alert_Context_Menu_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_context_menu_aria = /** @type {((inputs?: Alert_Context_Menu_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Context_Menu_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_context_menu_aria(inputs)
	return __es.alert_context_menu_aria(inputs)
});
/**
* | output |
* | --- |
* | "Created by {name}" |
*
* @param {Alert_Created_ByInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_created_by = /** @type {((inputs: Alert_Created_ByInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Created_ByInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_created_by(inputs)
	return __es.alert_created_by(inputs)
});
/**
* | output |
* | --- |
* | "Created through code" |
*
* @param {Alert_Created_Through_CodeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_created_through_code = /** @type {((inputs?: Alert_Created_Through_CodeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Created_Through_CodeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_created_through_code(inputs)
	return __es.alert_created_through_code(inputs)
});
/**
* | output |
* | --- |
* | "Criteria" |
*
* @param {Alert_CriteriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_criteria = /** @type {((inputs?: Alert_CriteriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_CriteriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_criteria(inputs)
	return __es.alert_criteria(inputs)
});
/**
* | output |
* | --- |
* | "Dashboard" |
*
* @param {Alert_DashboardInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_dashboard = /** @type {((inputs?: Alert_DashboardInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_DashboardInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_dashboard(inputs)
	return __es.alert_dashboard(inputs)
});
/**
* | output |
* | --- |
* | "Delete Alert" |
*
* @param {Alert_DeleteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_delete = /** @type {((inputs?: Alert_DeleteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_DeleteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_delete(inputs)
	return __es.alert_delete(inputs)
});
/**
* | output |
* | --- |
* | "Edit alert" |
*
* @param {Alert_EditInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_edit = /** @type {((inputs?: Alert_EditInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_EditInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_edit(inputs)
	return __es.alert_edit(inputs)
});
/**
* | output |
* | --- |
* | "Email notifications" |
*
* @param {Alert_Email_NotificationsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_email_notifications = /** @type {((inputs?: Alert_Email_NotificationsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Email_NotificationsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_email_notifications(inputs)
	return __es.alert_email_notifications(inputs)
});
/**
* | output |
* | --- |
* | "Filters ({count})" |
*
* @param {Alert_Filters_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_filters_label = /** @type {((inputs: Alert_Filters_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Filters_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_filters_label(inputs)
	return __es.alert_filters_label(inputs)
});
/**
* | output |
* | --- |
* | "Back" |
*
* @param {Alert_Form_BackInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_back = /** @type {((inputs?: Alert_Form_BackInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_BackInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_back(inputs)
	return __es.alert_form_back(inputs)
});
/**
* | output |
* | --- |
* | "Cancel" |
*
* @param {Alert_Form_CancelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_cancel = /** @type {((inputs?: Alert_Form_CancelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_CancelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_cancel(inputs)
	return __es.alert_form_cancel(inputs)
});
/**
* | output |
* | --- |
* | "Create" |
*
* @param {Alert_Form_CreateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_create = /** @type {((inputs?: Alert_Form_CreateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_CreateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_create(inputs)
	return __es.alert_form_create(inputs)
});
/**
* | output |
* | --- |
* | "Create Alert" |
*
* @param {Alert_Form_Create_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_create_title = /** @type {((inputs?: Alert_Form_Create_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Create_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_create_title(inputs)
	return __es.alert_form_create_title(inputs)
});
/**
* | output |
* | --- |
* | "Alert created" |
*
* @param {Alert_Form_CreatedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_created = /** @type {((inputs?: Alert_Form_CreatedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_CreatedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_created(inputs)
	return __es.alert_form_created(inputs)
});
/**
* | output |
* | --- |
* | "Trigger alert when these conditions are met" |
*
* @param {Alert_Form_Criteria_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_criteria_description = /** @type {((inputs?: Alert_Form_Criteria_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Criteria_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_criteria_description(inputs)
	return __es.alert_form_criteria_description(inputs)
});
/**
* | output |
* | --- |
* | "Criteria group operation" |
*
* @param {Alert_Form_Criteria_Group_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_criteria_group_aria = /** @type {((inputs?: Alert_Form_Criteria_Group_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Criteria_Group_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_criteria_group_aria(inputs)
	return __es.alert_form_criteria_group_aria(inputs)
});
/**
* | output |
* | --- |
* | "Criteria measure" |
*
* @param {Alert_Form_Criteria_Measure_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_criteria_measure_aria = /** @type {((inputs?: Alert_Form_Criteria_Measure_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Criteria_Measure_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_criteria_measure_aria(inputs)
	return __es.alert_form_criteria_measure_aria(inputs)
});
/**
* | output |
* | --- |
* | "Measure" |
*
* @param {Alert_Form_Criteria_Measure_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_criteria_measure_placeholder = /** @type {((inputs?: Alert_Form_Criteria_Measure_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Criteria_Measure_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_criteria_measure_placeholder(inputs)
	return __es.alert_form_criteria_measure_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Criteria operator" |
*
* @param {Alert_Form_Criteria_Operator_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_criteria_operator_aria = /** @type {((inputs?: Alert_Form_Criteria_Operator_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Criteria_Operator_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_criteria_operator_aria(inputs)
	return __es.alert_form_criteria_operator_aria(inputs)
});
/**
* | output |
* | --- |
* | "Operator" |
*
* @param {Alert_Form_Criteria_Operator_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_criteria_operator_placeholder = /** @type {((inputs?: Alert_Form_Criteria_Operator_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Criteria_Operator_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_criteria_operator_placeholder(inputs)
	return __es.alert_form_criteria_operator_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Alert Preview" |
*
* @param {Alert_Form_Criteria_Preview_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_criteria_preview_title = /** @type {((inputs?: Alert_Form_Criteria_Preview_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Criteria_Preview_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_criteria_preview_title(inputs)
	return __es.alert_form_criteria_preview_title(inputs)
});
/**
* | output |
* | --- |
* | "Criteria" |
*
* @param {Alert_Form_Criteria_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_criteria_title = /** @type {((inputs?: Alert_Form_Criteria_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Criteria_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_criteria_title(inputs)
	return __es.alert_form_criteria_title(inputs)
});
/**
* | output |
* | --- |
* | "Criteria type" |
*
* @param {Alert_Form_Criteria_Type_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_criteria_type_aria = /** @type {((inputs?: Alert_Form_Criteria_Type_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Criteria_Type_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_criteria_type_aria(inputs)
	return __es.alert_form_criteria_type_aria(inputs)
});
/**
* | output |
* | --- |
* | "type" |
*
* @param {Alert_Form_Criteria_Type_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_criteria_type_placeholder = /** @type {((inputs?: Alert_Form_Criteria_Type_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Criteria_Type_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_criteria_type_placeholder(inputs)
	return __es.alert_form_criteria_type_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Criteria value" |
*
* @param {Alert_Form_Criteria_Value_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_criteria_value_title = /** @type {((inputs?: Alert_Form_Criteria_Value_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Criteria_Value_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_criteria_value_title(inputs)
	return __es.alert_form_criteria_value_title(inputs)
});
/**
* | output |
* | --- |
* | "Filters" |
*
* @param {Alert_Form_Data_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_data_filters = /** @type {((inputs?: Alert_Form_Data_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Data_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_data_filters(inputs)
	return __es.alert_form_data_filters(inputs)
});
/**
* | output |
* | --- |
* | "Measure" |
*
* @param {Alert_Form_Data_MeasureInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_data_measure = /** @type {((inputs?: Alert_Form_Data_MeasureInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Data_MeasureInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_data_measure(inputs)
	return __es.alert_form_data_measure(inputs)
});
/**
* | output |
* | --- |
* | "Select a measure" |
*
* @param {Alert_Form_Data_Measure_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_data_measure_placeholder = /** @type {((inputs?: Alert_Form_Data_Measure_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Data_Measure_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_data_measure_placeholder(inputs)
	return __es.alert_form_data_measure_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Select the measures you want to monitor." |
*
* @param {Alert_Form_Data_Measures_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_data_measures_desc = /** @type {((inputs?: Alert_Form_Data_Measures_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Data_Measures_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_data_measures_desc(inputs)
	return __es.alert_form_data_measures_desc(inputs)
});
/**
* | output |
* | --- |
* | "None" |
*
* @param {Alert_Form_Data_NoneInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_data_none = /** @type {((inputs?: Alert_Form_Data_NoneInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Data_NoneInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_data_none(inputs)
	return __es.alert_form_data_none(inputs)
});
/**
* | output |
* | --- |
* | "Data preview" |
*
* @param {Alert_Form_Data_PreviewInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_data_preview = /** @type {((inputs?: Alert_Form_Data_PreviewInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Data_PreviewInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_data_preview(inputs)
	return __es.alert_form_data_preview(inputs)
});
/**
* | output |
* | --- |
* | "Here's a look at the data you've selected above." |
*
* @param {Alert_Form_Data_Preview_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_data_preview_desc = /** @type {((inputs?: Alert_Form_Data_Preview_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Data_Preview_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_data_preview_desc(inputs)
	return __es.alert_form_data_preview_desc(inputs)
});
/**
* | output |
* | --- |
* | "Split by dimension" |
*
* @param {Alert_Form_Data_Split_ByInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_data_split_by = /** @type {((inputs?: Alert_Form_Data_Split_ByInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Data_Split_ByInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_data_split_by(inputs)
	return __es.alert_form_data_split_by(inputs)
});
/**
* | output |
* | --- |
* | "Select a dimension" |
*
* @param {Alert_Form_Data_Split_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_data_split_placeholder = /** @type {((inputs?: Alert_Form_Data_Split_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Data_Split_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_data_split_placeholder(inputs)
	return __es.alert_form_data_split_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Alert data" |
*
* @param {Alert_Form_Data_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_data_title = /** @type {((inputs?: Alert_Form_Data_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Data_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_data_title(inputs)
	return __es.alert_form_data_title(inputs)
});
/**
* | output |
* | --- |
* | "Edit Alert" |
*
* @param {Alert_Form_Edit_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_edit_title = /** @type {((inputs?: Alert_Form_Edit_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Edit_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_edit_title(inputs)
	return __es.alert_form_edit_title(inputs)
});
/**
* | output |
* | --- |
* | "Alert edited" |
*
* @param {Alert_Form_EditedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_edited = /** @type {((inputs?: Alert_Form_EditedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_EditedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_edited(inputs)
	return __es.alert_form_edited(inputs)
});
/**
* | output |
* | --- |
* | "We'll email alerts to these addresses. Make sure they have access to your project." |
*
* @param {Alert_Form_Email_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_email_desc = /** @type {((inputs?: Alert_Form_Email_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Email_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_email_desc(inputs)
	return __es.alert_form_email_desc(inputs)
});
/**
* | output |
* | --- |
* | "Enter an email address" |
*
* @param {Alert_Form_Email_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_email_placeholder = /** @type {((inputs?: Alert_Form_Email_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Email_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_email_placeholder(inputs)
	return __es.alert_form_email_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Email notifications" |
*
* @param {Alert_Form_Email_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_email_title = /** @type {((inputs?: Alert_Form_Email_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Email_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_email_title(inputs)
	return __es.alert_form_email_title(inputs)
});
/**
* | output |
* | --- |
* | "Go to alerts" |
*
* @param {Alert_Form_Go_To_AlertsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_go_to_alerts = /** @type {((inputs?: Alert_Form_Go_To_AlertsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Go_To_AlertsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_go_to_alerts(inputs)
	return __es.alert_form_go_to_alerts(inputs)
});
/**
* | output |
* | --- |
* | "My alert" |
*
* @param {Alert_Form_Name_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_name_placeholder = /** @type {((inputs?: Alert_Form_Name_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Name_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_name_placeholder(inputs)
	return __es.alert_form_name_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Alert name" |
*
* @param {Alert_Form_Name_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_name_title = /** @type {((inputs?: Alert_Form_Name_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Name_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_name_title(inputs)
	return __es.alert_form_name_title(inputs)
});
/**
* | output |
* | --- |
* | "Next" |
*
* @param {Alert_Form_NextInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_next = /** @type {((inputs?: Alert_Form_NextInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_NextInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_next(inputs)
	return __es.alert_form_next(inputs)
});
/**
* | output |
* | --- |
* | "No criteria selected" |
*
* @param {Alert_Form_No_CriteriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_no_criteria = /** @type {((inputs?: Alert_Form_No_CriteriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_No_CriteriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_no_criteria(inputs)
	return __es.alert_form_no_criteria(inputs)
});
/**
* | output |
* | --- |
* | "No data to preview" |
*
* @param {Alert_Form_No_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_no_data = /** @type {((inputs?: Alert_Form_No_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_No_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_no_data(inputs)
	return __es.alert_form_no_data(inputs)
});
/**
* | output |
* | --- |
* | "Preview cell" |
*
* @param {Alert_Form_Preview_CellInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_preview_cell = /** @type {((inputs?: Alert_Form_Preview_CellInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Preview_CellInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_preview_cell(inputs)
	return __es.alert_form_preview_cell(inputs)
});
/**
* | output |
* | --- |
* | "alert preview table" |
*
* @param {Alert_Form_Preview_Table_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_preview_table_aria = /** @type {((inputs?: Alert_Form_Preview_Table_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Preview_Table_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_preview_table_aria(inputs)
	return __es.alert_form_preview_table_aria(inputs)
});
/**
* | output |
* | --- |
* | "Select criteria to see a preview" |
*
* @param {Alert_Form_Select_CriteriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_select_criteria = /** @type {((inputs?: Alert_Form_Select_CriteriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Select_CriteriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_select_criteria(inputs)
	return __es.alert_form_select_criteria(inputs)
});
/**
* | output |
* | --- |
* | "We'll send alerts directly to these channels." |
*
* @param {Alert_Form_Slack_Channels_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_slack_channels_desc = /** @type {((inputs?: Alert_Form_Slack_Channels_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Slack_Channels_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_slack_channels_desc(inputs)
	return __es.alert_form_slack_channels_desc(inputs)
});
/**
* | output |
* | --- |
* | "Slack has not been configured for this project. Read the <a href=\"{docsUrl}\" target=\"_blank\">docs</a> to learn more." |
*
* @param {Alert_Form_Slack_Not_ConfiguredInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_slack_not_configured = /** @type {((inputs: Alert_Form_Slack_Not_ConfiguredInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Slack_Not_ConfiguredInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_slack_not_configured(inputs)
	return __es.alert_form_slack_not_configured(inputs)
});
/**
* | output |
* | --- |
* | "# Enter a Slack channel name" |
*
* @param {Alert_Form_Slack_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_slack_placeholder = /** @type {((inputs?: Alert_Form_Slack_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Slack_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_slack_placeholder(inputs)
	return __es.alert_form_slack_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Slack notifications" |
*
* @param {Alert_Form_Slack_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_slack_title = /** @type {((inputs?: Alert_Form_Slack_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Slack_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_slack_title(inputs)
	return __es.alert_form_slack_title(inputs)
});
/**
* | output |
* | --- |
* | "We'll alert them with direct messages in Slack." |
*
* @param {Alert_Form_Slack_Users_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_slack_users_desc = /** @type {((inputs?: Alert_Form_Slack_Users_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Slack_Users_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_slack_users_desc(inputs)
	return __es.alert_form_slack_users_desc(inputs)
});
/**
* | output |
* | --- |
* | "Set a snooze period to silence repeat notifications for the same alert." |
*
* @param {Alert_Form_Snooze_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_snooze_desc = /** @type {((inputs?: Alert_Form_Snooze_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Snooze_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_snooze_desc(inputs)
	return __es.alert_form_snooze_desc(inputs)
});
/**
* | output |
* | --- |
* | "Snooze" |
*
* @param {Alert_Form_Snooze_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_snooze_title = /** @type {((inputs?: Alert_Form_Snooze_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Snooze_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_snooze_title(inputs)
	return __es.alert_form_snooze_title(inputs)
});
/**
* | output |
* | --- |
* | "Criteria" |
*
* @param {Alert_Form_Tab_CriteriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_tab_criteria = /** @type {((inputs?: Alert_Form_Tab_CriteriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Tab_CriteriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_tab_criteria(inputs)
	return __es.alert_form_tab_criteria(inputs)
});
/**
* | output |
* | --- |
* | "Data" |
*
* @param {Alert_Form_Tab_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_tab_data = /** @type {((inputs?: Alert_Form_Tab_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Tab_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_tab_data(inputs)
	return __es.alert_form_tab_data(inputs)
});
/**
* | output |
* | --- |
* | "Delivery" |
*
* @param {Alert_Form_Tab_DeliveryInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_tab_delivery = /** @type {((inputs?: Alert_Form_Tab_DeliveryInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Tab_DeliveryInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_tab_delivery(inputs)
	return __es.alert_form_tab_delivery(inputs)
});
/**
* | output |
* | --- |
* | "Trigger" |
*
* @param {Alert_Form_TriggerInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_trigger = /** @type {((inputs?: Alert_Form_TriggerInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_TriggerInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_trigger(inputs)
	return __es.alert_form_trigger(inputs)
});
/**
* | output |
* | --- |
* | "Whenever data refreshes" |
*
* @param {Alert_Form_Trigger_Data_RefreshInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_trigger_data_refresh = /** @type {((inputs?: Alert_Form_Trigger_Data_RefreshInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Trigger_Data_RefreshInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_trigger_data_refresh(inputs)
	return __es.alert_form_trigger_data_refresh(inputs)
});
/**
* | output |
* | --- |
* | "Set schedule" |
*
* @param {Alert_Form_Trigger_Set_ScheduleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_trigger_set_schedule = /** @type {((inputs?: Alert_Form_Trigger_Set_ScheduleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_Trigger_Set_ScheduleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_trigger_set_schedule(inputs)
	return __es.alert_form_trigger_set_schedule(inputs)
});
/**
* | output |
* | --- |
* | "Update" |
*
* @param {Alert_Form_UpdateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_form_update = /** @type {((inputs?: Alert_Form_UpdateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Form_UpdateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_form_update(inputs)
	return __es.alert_form_update(inputs)
});
/**
* | output |
* | --- |
* | "Last checked {time}" |
*
* @param {Alert_Last_CheckedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_last_checked = /** @type {((inputs: Alert_Last_CheckedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Last_CheckedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_last_checked(inputs)
	return __es.alert_last_checked(inputs)
});
/**
* | output |
* | --- |
* | "Name" |
*
* @param {Alert_Name_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_name_label = /** @type {((inputs?: Alert_Name_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Name_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_name_label(inputs)
	return __es.alert_name_label(inputs)
});
/**
* | output |
* | --- |
* | "No filters were applied to the dashboard that this alert was created from." |
*
* @param {Alert_No_Filters_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_no_filters_body = /** @type {((inputs?: Alert_No_Filters_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_No_Filters_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_no_filters_body(inputs)
	return __es.alert_no_filters_body(inputs)
});
/**
* | output |
* | --- |
* | "No filters selected" |
*
* @param {Alert_No_Filters_HeadingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_no_filters_heading = /** @type {((inputs?: Alert_No_Filters_HeadingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_No_Filters_HeadingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_no_filters_heading(inputs)
	return __es.alert_no_filters_heading(inputs)
});
/**
* | output |
* | --- |
* | "To apply filters, close this window and filter your dashboard." |
*
* @param {Alert_No_Filters_HintInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_no_filters_hint = /** @type {((inputs?: Alert_No_Filters_HintInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_No_Filters_HintInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_no_filters_hint(inputs)
	return __es.alert_no_filters_hint(inputs)
});
/**
* | output |
* | --- |
* | "None" |
*
* @param {Alert_NoneInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_none = /** @type {((inputs?: Alert_NoneInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_NoneInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_none(inputs)
	return __es.alert_none(inputs)
});
/**
* | output |
* | --- |
* | "Hasn't been checked yet" |
*
* @param {Alert_Not_Checked_YetInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_not_checked_yet = /** @type {((inputs?: Alert_Not_Checked_YetInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Not_Checked_YetInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_not_checked_yet(inputs)
	return __es.alert_not_checked_yet(inputs)
});
/**
* | output |
* | --- |
* | "Schedule" |
*
* @param {Alert_ScheduleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_schedule = /** @type {((inputs?: Alert_ScheduleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_ScheduleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_schedule(inputs)
	return __es.alert_schedule(inputs)
});
/**
* | output |
* | --- |
* | "Slack notifications" |
*
* @param {Alert_Slack_NotificationsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_slack_notifications = /** @type {((inputs?: Alert_Slack_NotificationsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Slack_NotificationsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_slack_notifications(inputs)
	return __es.alert_slack_notifications(inputs)
});
/**
* | output |
* | --- |
* | "Snooze" |
*
* @param {Alert_SnoozeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_snooze = /** @type {((inputs?: Alert_SnoozeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_SnoozeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_snooze(inputs)
	return __es.alert_snooze(inputs)
});
/**
* | output |
* | --- |
* | "Split by dimension" |
*
* @param {Alert_Split_By_DimensionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_split_by_dimension = /** @type {((inputs?: Alert_Split_By_DimensionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Split_By_DimensionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_split_by_dimension(inputs)
	return __es.alert_split_by_dimension(inputs)
});
/**
* | output |
* | --- |
* | "Checked" |
*
* @param {Alert_Status_CheckedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_status_checked = /** @type {((inputs?: Alert_Status_CheckedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Status_CheckedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_status_checked(inputs)
	return __es.alert_status_checked(inputs)
});
/**
* | output |
* | --- |
* | "Checking" |
*
* @param {Alert_Status_CheckingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_status_checking = /** @type {((inputs?: Alert_Status_CheckingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Status_CheckingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_status_checking(inputs)
	return __es.alert_status_checking(inputs)
});
/**
* | output |
* | --- |
* | "Failed" |
*
* @param {Alert_Status_FailedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_status_failed = /** @type {((inputs?: Alert_Status_FailedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Status_FailedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_status_failed(inputs)
	return __es.alert_status_failed(inputs)
});
/**
* | output |
* | --- |
* | "Not triggered" |
*
* @param {Alert_Status_Not_TriggeredInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_status_not_triggered = /** @type {((inputs?: Alert_Status_Not_TriggeredInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Status_Not_TriggeredInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_status_not_triggered(inputs)
	return __es.alert_status_not_triggered(inputs)
});
/**
* | output |
* | --- |
* | "Running" |
*
* @param {Alert_Status_RunningInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_status_running = /** @type {((inputs?: Alert_Status_RunningInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Status_RunningInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_status_running(inputs)
	return __es.alert_status_running(inputs)
});
/**
* | output |
* | --- |
* | "Triggered" |
*
* @param {Alert_Status_TriggeredInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_status_triggered = /** @type {((inputs?: Alert_Status_TriggeredInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Status_TriggeredInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_status_triggered(inputs)
	return __es.alert_status_triggered(inputs)
});
/**
* | output |
* | --- |
* | "Status unknown" |
*
* @param {Alert_Status_UnknownInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_status_unknown = /** @type {((inputs?: Alert_Status_UnknownInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Status_UnknownInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_status_unknown(inputs)
	return __es.alert_status_unknown(inputs)
});
/**
* | output |
* | --- |
* | "Failed to unsubscribe." |
*
* @param {Alert_Unsubscribe_FailedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_unsubscribe_failed = /** @type {((inputs?: Alert_Unsubscribe_FailedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Unsubscribe_FailedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_unsubscribe_failed(inputs)
	return __es.alert_unsubscribe_failed(inputs)
});
/**
* | output |
* | --- |
* | "Unsubscribed from alert." |
*
* @param {Alert_UnsubscribedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_unsubscribed = /** @type {((inputs?: Alert_UnsubscribedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_UnsubscribedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_unsubscribed(inputs)
	return __es.alert_unsubscribed(inputs)
});
/**
* | output |
* | --- |
* | "Unsubscribing..." |
*
* @param {Alert_UnsubscribingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_unsubscribing = /** @type {((inputs?: Alert_UnsubscribingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_UnsubscribingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_unsubscribing(inputs)
	return __es.alert_unsubscribing(inputs)
});
/**
* | output |
* | --- |
* | "Whenever your data refreshes" |
*
* @param {Alert_Whenever_Data_RefreshesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alert_whenever_data_refreshes = /** @type {((inputs?: Alert_Whenever_Data_RefreshesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alert_Whenever_Data_RefreshesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alert_whenever_data_refreshes(inputs)
	return __es.alert_whenever_data_refreshes(inputs)
});
/**
* | output |
* | --- |
* | "Create {alertsLink} from any dashboard or {codeLink}." |
*
* @param {Alerts_Empty_ActionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alerts_empty_action = /** @type {((inputs: Alerts_Empty_ActionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alerts_Empty_ActionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alerts_empty_action(inputs)
	return __es.alerts_empty_action(inputs)
});
/**
* | output |
* | --- |
* | "You don't have any alerts yet" |
*
* @param {Alerts_Empty_MessageInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alerts_empty_message = /** @type {((inputs?: Alerts_Empty_MessageInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alerts_Empty_MessageInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alerts_empty_message(inputs)
	return __es.alerts_empty_message(inputs)
});
/**
* | output |
* | --- |
* | "alerts" |
*
* @param {Alerts_Link_TextInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alerts_link_text = /** @type {((inputs?: Alerts_Link_TextInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alerts_Link_TextInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alerts_link_text(inputs)
	return __es.alerts_link_text(inputs)
});
/**
* | output |
* | --- |
* | "via code" |
*
* @param {Alerts_Via_CodeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const alerts_via_code = /** @type {((inputs?: Alerts_Via_CodeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Alerts_Via_CodeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.alerts_via_code(inputs)
	return __es.alerts_via_code(inputs)
});
/**
* | output |
* | --- |
* | "Contact Rill support" |
*
* @param {Avatar_Contact_SupportInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const avatar_contact_support = /** @type {((inputs?: Avatar_Contact_SupportInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Avatar_Contact_SupportInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.avatar_contact_support(inputs)
	return __es.avatar_contact_support(inputs)
});
/**
* | output |
* | --- |
* | "Copied URL" |
*
* @param {Avatar_Copied_UrlInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const avatar_copied_url = /** @type {((inputs?: Avatar_Copied_UrlInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Avatar_Copied_UrlInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.avatar_copied_url(inputs)
	return __es.avatar_copied_url(inputs)
});
/**
* | output |
* | --- |
* | "Copy URL" |
*
* @param {Avatar_Copy_UrlInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const avatar_copy_url = /** @type {((inputs?: Avatar_Copy_UrlInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Avatar_Copy_UrlInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.avatar_copy_url(inputs)
	return __es.avatar_copy_url(inputs)
});
/**
* | output |
* | --- |
* | "Copy URL for this view" |
*
* @param {Avatar_Copy_Url_For_ViewInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const avatar_copy_url_for_view = /** @type {((inputs?: Avatar_Copy_Url_For_ViewInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Avatar_Copy_Url_For_ViewInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.avatar_copy_url_for_view(inputs)
	return __es.avatar_copy_url_for_view(inputs)
});
/**
* | output |
* | --- |
* | "Create public URL" |
*
* @param {Avatar_Create_Public_UrlInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const avatar_create_public_url = /** @type {((inputs?: Avatar_Create_Public_UrlInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Avatar_Create_Public_UrlInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.avatar_create_public_url(inputs)
	return __es.avatar_create_public_url(inputs)
});
/**
* | output |
* | --- |
* | "Documentation" |
*
* @param {Avatar_DocumentationInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const avatar_documentation = /** @type {((inputs?: Avatar_DocumentationInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Avatar_DocumentationInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.avatar_documentation(inputs)
	return __es.avatar_documentation(inputs)
});
/**
* | output |
* | --- |
* | "Join us on Discord" |
*
* @param {Avatar_Join_DiscordInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const avatar_join_discord = /** @type {((inputs?: Avatar_Join_DiscordInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Avatar_Join_DiscordInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.avatar_join_discord(inputs)
	return __es.avatar_join_discord(inputs)
});
/**
* | output |
* | --- |
* | "Logout" |
*
* @param {Avatar_LogoutInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const avatar_logout = /** @type {((inputs?: Avatar_LogoutInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Avatar_LogoutInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.avatar_logout(inputs)
	return __es.avatar_logout(inputs)
});
/**
* | output |
* | --- |
* | "Share" |
*
* @param {Avatar_ShareInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const avatar_share = /** @type {((inputs?: Avatar_ShareInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Avatar_ShareInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.avatar_share(inputs)
	return __es.avatar_share(inputs)
});
/**
* | output |
* | --- |
* | "Share dashboard" |
*
* @param {Avatar_Share_DashboardInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const avatar_share_dashboard = /** @type {((inputs?: Avatar_Share_DashboardInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Avatar_Share_DashboardInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.avatar_share_dashboard(inputs)
	return __es.avatar_share_dashboard(inputs)
});
/**
* | output |
* | --- |
* | "Share your current view with another project member." |
*
* @param {Avatar_Share_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const avatar_share_description = /** @type {((inputs?: Avatar_Share_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Avatar_Share_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.avatar_share_description(inputs)
	return __es.avatar_share_description(inputs)
});
/**
* | output |
* | --- |
* | "View as" |
*
* @param {Avatar_View_AsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const avatar_view_as = /** @type {((inputs?: Avatar_View_AsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Avatar_View_AsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.avatar_view_as(inputs)
	return __es.avatar_view_as(inputs)
});
/**
* | output |
* | --- |
* | "Copy this value to clipboard" |
*
* @param {Bignumber_Copy_ValueInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bignumber_copy_value = /** @type {((inputs?: Bignumber_Copy_ValueInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bignumber_Copy_ValueInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bignumber_copy_value(inputs)
	return __es.bignumber_copy_value(inputs)
});
/**
* | output |
* | --- |
* | "+ Click" |
*
* @param {Bignumber_Shift_ClickInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bignumber_shift_click = /** @type {((inputs?: Bignumber_Shift_ClickInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bignumber_Shift_ClickInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bignumber_shift_click(inputs)
	return __es.bignumber_shift_click(inputs)
});
/**
* | output |
* | --- |
* | "Absolute time range" |
*
* @param {Bookmark_Absolute_Time_RangeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_absolute_time_range = /** @type {((inputs?: Bookmark_Absolute_Time_RangeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Absolute_Time_RangeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_absolute_time_range(inputs)
	return __es.bookmark_absolute_time_range(inputs)
});
/**
* | output |
* | --- |
* | "The bookmark will use the dashboard's relative time if this toggle is off." |
*
* @param {Bookmark_Absolute_Time_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_absolute_time_tooltip = /** @type {((inputs?: Bookmark_Absolute_Time_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Absolute_Time_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_absolute_time_tooltip(inputs)
	return __es.bookmark_absolute_time_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "Category" |
*
* @param {Bookmark_CategoryInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_category = /** @type {((inputs?: Bookmark_CategoryInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_CategoryInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_category(inputs)
	return __es.bookmark_category(inputs)
});
/**
* | output |
* | --- |
* | "Your bookmarks can only be viewed by you. Managed bookmarks will be available to all viewers of this dashboard." |
*
* @param {Bookmark_Category_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_category_tooltip = /** @type {((inputs?: Bookmark_Category_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Category_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_category_tooltip(inputs)
	return __es.bookmark_category_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "Bookmark created" |
*
* @param {Bookmark_CreatedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_created = /** @type {((inputs?: Bookmark_CreatedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_CreatedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_created(inputs)
	return __es.bookmark_created(inputs)
});
/**
* | output |
* | --- |
* | "Created by project admin" |
*
* @param {Bookmark_Created_By_AdminInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_created_by_admin = /** @type {((inputs?: Bookmark_Created_By_AdminInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Created_By_AdminInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_created_by_admin(inputs)
	return __es.bookmark_created_by_admin(inputs)
});
/**
* | output |
* | --- |
* | "Bookmark current view" |
*
* @param {Bookmark_Current_ViewInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_current_view = /** @type {((inputs?: Bookmark_Current_ViewInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Current_ViewInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_current_view(inputs)
	return __es.bookmark_current_view(inputs)
});
/**
* | output |
* | --- |
* | "Bookmark current view as Home." |
*
* @param {Bookmark_Current_View_As_HomeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_current_view_as_home = /** @type {((inputs?: Bookmark_Current_View_As_HomeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Current_View_As_HomeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_current_view_as_home(inputs)
	return __es.bookmark_current_view_as_home(inputs)
});
/**
* | output |
* | --- |
* | "Default Label" |
*
* @param {Bookmark_Default_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_default_label = /** @type {((inputs?: Bookmark_Default_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Default_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_default_label(inputs)
	return __es.bookmark_default_label(inputs)
});
/**
* | output |
* | --- |
* | "Delete bookmark" |
*
* @param {Bookmark_Delete_BookmarkInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_delete_bookmark = /** @type {((inputs?: Bookmark_Delete_BookmarkInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Delete_BookmarkInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_delete_bookmark(inputs)
	return __es.bookmark_delete_bookmark(inputs)
});
/**
* | output |
* | --- |
* | "Delete Home bookmark" |
*
* @param {Bookmark_Delete_Home_BookmarkInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_delete_home_bookmark = /** @type {((inputs?: Bookmark_Delete_Home_BookmarkInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Delete_Home_BookmarkInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_delete_home_bookmark(inputs)
	return __es.bookmark_delete_home_bookmark(inputs)
});
/**
* | output |
* | --- |
* | "Bookmark {name} deleted" |
*
* @param {Bookmark_DeletedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_deleted = /** @type {((inputs: Bookmark_DeletedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_DeletedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_deleted(inputs)
	return __es.bookmark_deleted(inputs)
});
/**
* | output |
* | --- |
* | "Description" |
*
* @param {Bookmark_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_description = /** @type {((inputs?: Bookmark_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_description(inputs)
	return __es.bookmark_description(inputs)
});
/**
* | output |
* | --- |
* | "Edit bookmark" |
*
* @param {Bookmark_EditInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_edit = /** @type {((inputs?: Bookmark_EditInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_EditInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_edit(inputs)
	return __es.bookmark_edit(inputs)
});
/**
* | output |
* | --- |
* | "Filters" |
*
* @param {Bookmark_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_filters = /** @type {((inputs?: Bookmark_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_filters(inputs)
	return __es.bookmark_filters(inputs)
});
/**
* | output |
* | --- |
* | "Inherited from underlying dashboard view." |
*
* @param {Bookmark_Filters_InheritedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_filters_inherited = /** @type {((inputs?: Bookmark_Filters_InheritedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Filters_InheritedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_filters_inherited(inputs)
	return __es.bookmark_filters_inherited(inputs)
});
/**
* | output |
* | --- |
* | "Toggling this on will only save the filter set above, not the full dashboard layout and state." |
*
* @param {Bookmark_Filters_Only_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_filters_only_tooltip = /** @type {((inputs?: Bookmark_Filters_Only_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Filters_Only_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_filters_only_tooltip(inputs)
	return __es.bookmark_filters_only_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "Go to Home" |
*
* @param {Bookmark_Go_To_HomeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_go_to_home = /** @type {((inputs?: Bookmark_Go_To_HomeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Go_To_HomeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_go_to_home(inputs)
	return __es.bookmark_go_to_home(inputs)
});
/**
* | output |
* | --- |
* | "Home bookmark created" |
*
* @param {Bookmark_Home_CreatedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_home_created = /** @type {((inputs?: Bookmark_Home_CreatedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Home_CreatedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_home_created(inputs)
	return __es.bookmark_home_created(inputs)
});
/**
* | output |
* | --- |
* | "This will be everyone's main view for this dashboard." |
*
* @param {Bookmark_Home_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_home_description = /** @type {((inputs?: Bookmark_Home_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Home_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_home_description(inputs)
	return __es.bookmark_home_description(inputs)
});
/**
* | output |
* | --- |
* | "Label" |
*
* @param {Bookmark_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_label = /** @type {((inputs?: Bookmark_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_label(inputs)
	return __es.bookmark_label(inputs)
});
/**
* | output |
* | --- |
* | "Managed bookmarks" |
*
* @param {Bookmark_Managed_BookmarksInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_managed_bookmarks = /** @type {((inputs?: Bookmark_Managed_BookmarksInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Managed_BookmarksInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_managed_bookmarks(inputs)
	return __es.bookmark_managed_bookmarks(inputs)
});
/**
* | output |
* | --- |
* | "You have no bookmarks for this dashboard." |
*
* @param {Bookmark_No_BookmarksInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_no_bookmarks = /** @type {((inputs?: Bookmark_No_BookmarksInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_No_BookmarksInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_no_bookmarks(inputs)
	return __es.bookmark_no_bookmarks(inputs)
});
/**
* | output |
* | --- |
* | "There are no shared bookmarks for this dashboard." |
*
* @param {Bookmark_No_SharedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_no_shared = /** @type {((inputs?: Bookmark_No_SharedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_No_SharedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_no_shared(inputs)
	return __es.bookmark_no_shared(inputs)
});
/**
* | output |
* | --- |
* | "Return to dashboard home" |
*
* @param {Bookmark_Return_To_HomeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_return_to_home = /** @type {((inputs?: Bookmark_Return_To_HomeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Return_To_HomeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_return_to_home(inputs)
	return __es.bookmark_return_to_home(inputs)
});
/**
* | output |
* | --- |
* | "Save" |
*
* @param {Bookmark_SaveInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_save = /** @type {((inputs?: Bookmark_SaveInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_SaveInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_save(inputs)
	return __es.bookmark_save(inputs)
});
/**
* | output |
* | --- |
* | "Save filters only" |
*
* @param {Bookmark_Save_Filters_OnlyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_save_filters_only = /** @type {((inputs?: Bookmark_Save_Filters_OnlyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Save_Filters_OnlyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_save_filters_only(inputs)
	return __es.bookmark_save_filters_only(inputs)
});
/**
* | output |
* | --- |
* | "Bookmark updated" |
*
* @param {Bookmark_UpdatedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_updated = /** @type {((inputs?: Bookmark_UpdatedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_UpdatedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_updated(inputs)
	return __es.bookmark_updated(inputs)
});
/**
* | output |
* | --- |
* | "Your bookmarks" |
*
* @param {Bookmark_Your_BookmarksInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const bookmark_your_bookmarks = /** @type {((inputs?: Bookmark_Your_BookmarksInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Bookmark_Your_BookmarksInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.bookmark_your_bookmarks(inputs)
	return __es.bookmark_your_bookmarks(inputs)
});
/**
* | output |
* | --- |
* | "Apply" |
*
* @param {Calendar_ApplyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const calendar_apply = /** @type {((inputs?: Calendar_ApplyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Calendar_ApplyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.calendar_apply(inputs)
	return __es.calendar_apply(inputs)
});
/**
* | output |
* | --- |
* | "Range exceeds the {capLabel} query limit." |
*
* @param {Calendar_Range_Exceeds_LimitInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const calendar_range_exceeds_limit = /** @type {((inputs: Calendar_Range_Exceeds_LimitInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Calendar_Range_Exceeds_LimitInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.calendar_range_exceeds_limit(inputs)
	return __es.calendar_range_exceeds_limit(inputs)
});
/**
* | output |
* | --- |
* | "Add widget" |
*
* @param {Canvas_Add_WidgetInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_add_widget = /** @type {((inputs?: Canvas_Add_WidgetInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Add_WidgetInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_add_widget(inputs)
	return __es.canvas_add_widget(inputs)
});
/**
* | output |
* | --- |
* | "AI is generating chart" |
*
* @param {Canvas_Ai_Generating_ChartInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_ai_generating_chart = /** @type {((inputs?: Canvas_Ai_Generating_ChartInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Ai_Generating_ChartInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_ai_generating_chart(inputs)
	return __es.canvas_ai_generating_chart(inputs)
});
/**
* | output |
* | --- |
* | "Opens the AI assistant to edit this chart" |
*
* @param {Canvas_Ai_HintInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_ai_hint = /** @type {((inputs?: Canvas_Ai_HintInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Ai_HintInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_ai_hint(inputs)
	return __es.canvas_ai_hint(inputs)
});
/**
* | output |
* | --- |
* | "AI is editing" |
*
* @param {Canvas_Ai_Is_EditingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_ai_is_editing = /** @type {((inputs?: Canvas_Ai_Is_EditingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Ai_Is_EditingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_ai_is_editing(inputs)
	return __es.canvas_ai_is_editing(inputs)
});
/**
* | output |
* | --- |
* | "Use the inspector panel to write Metrics SQL and Vega-Lite spec manually." |
*
* @param {Canvas_Ai_Manual_HintInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_ai_manual_hint = /** @type {((inputs?: Canvas_Ai_Manual_HintInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Ai_Manual_HintInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_ai_manual_hint(inputs)
	return __es.canvas_ai_manual_hint(inputs)
});
/**
* | output |
* | --- |
* | "Write SQL & Vega-Lite manually" |
*
* @param {Canvas_Ai_Write_ManuallyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_ai_write_manually = /** @type {((inputs?: Canvas_Ai_Write_ManuallyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Ai_Write_ManuallyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_ai_write_manually(inputs)
	return __es.canvas_ai_write_manually(inputs)
});
/**
* | output |
* | --- |
* | "Align bottom" |
*
* @param {Canvas_Align_BottomInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_align_bottom = /** @type {((inputs?: Canvas_Align_BottomInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Align_BottomInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_align_bottom(inputs)
	return __es.canvas_align_bottom(inputs)
});
/**
* | output |
* | --- |
* | "Align center" |
*
* @param {Canvas_Align_CenterInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_align_center = /** @type {((inputs?: Canvas_Align_CenterInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Align_CenterInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_align_center(inputs)
	return __es.canvas_align_center(inputs)
});
/**
* | output |
* | --- |
* | "Align left" |
*
* @param {Canvas_Align_LeftInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_align_left = /** @type {((inputs?: Canvas_Align_LeftInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Align_LeftInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_align_left(inputs)
	return __es.canvas_align_left(inputs)
});
/**
* | output |
* | --- |
* | "Align middle" |
*
* @param {Canvas_Align_MiddleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_align_middle = /** @type {((inputs?: Canvas_Align_MiddleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Align_MiddleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_align_middle(inputs)
	return __es.canvas_align_middle(inputs)
});
/**
* | output |
* | --- |
* | "Align right" |
*
* @param {Canvas_Align_RightInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_align_right = /** @type {((inputs?: Canvas_Align_RightInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Align_RightInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_align_right(inputs)
	return __es.canvas_align_right(inputs)
});
/**
* | output |
* | --- |
* | "Align top" |
*
* @param {Canvas_Align_TopInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_align_top = /** @type {((inputs?: Canvas_Align_TopInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Align_TopInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_align_top(inputs)
	return __es.canvas_align_top(inputs)
});
/**
* | output |
* | --- |
* | "Alignment" |
*
* @param {Canvas_Alignment_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_alignment_label = /** @type {((inputs?: Canvas_Alignment_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Alignment_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_alignment_label(inputs)
	return __es.canvas_alignment_label(inputs)
});
/**
* | output |
* | --- |
* | "Apply measure value formatting" |
*
* @param {Canvas_Apply_Measure_Formatting_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_apply_measure_formatting_label = /** @type {((inputs?: Canvas_Apply_Measure_Formatting_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Apply_Measure_Formatting_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_apply_measure_formatting_label(inputs)
	return __es.canvas_apply_measure_formatting_label(inputs)
});
/**
* | output |
* | --- |
* | "← Back to prompt" |
*
* @param {Canvas_Back_To_PromptInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_back_to_prompt = /** @type {((inputs?: Canvas_Back_To_PromptInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Back_To_PromptInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_back_to_prompt(inputs)
	return __es.canvas_back_to_prompt(inputs)
});
/**
* | output |
* | --- |
* | "Blues" |
*
* @param {Canvas_BluesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_blues = /** @type {((inputs?: Canvas_BluesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_BluesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_blues(inputs)
	return __es.canvas_blues(inputs)
});
/**
* | output |
* | --- |
* | "Breakdown by" |
*
* @param {Canvas_Breakdown_By_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_breakdown_by_label = /** @type {((inputs?: Canvas_Breakdown_By_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Breakdown_By_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_breakdown_by_label(inputs)
	return __es.canvas_breakdown_by_label(inputs)
});
/**
* | output |
* | --- |
* | "Cancel" |
*
* @param {Canvas_CancelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_cancel = /** @type {((inputs?: Canvas_CancelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_CancelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_cancel(inputs)
	return __es.canvas_cancel(inputs)
});
/**
* | output |
* | --- |
* | "Chart" |
*
* @param {Canvas_ChartInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_chart = /** @type {((inputs?: Canvas_ChartInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_ChartInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_chart(inputs)
	return __es.canvas_chart(inputs)
});
/**
* | output |
* | --- |
* | "Chart type" |
*
* @param {Canvas_Chart_TypeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_chart_type = /** @type {((inputs?: Canvas_Chart_TypeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Chart_TypeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_chart_type(inputs)
	return __es.canvas_chart_type(inputs)
});
/**
* | output |
* | --- |
* | "Cividis" |
*
* @param {Canvas_CividisInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_cividis = /** @type {((inputs?: Canvas_CividisInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_CividisInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_cividis(inputs)
	return __es.canvas_cividis(inputs)
});
/**
* | output |
* | --- |
* | "Clear filters" |
*
* @param {Canvas_Clear_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_clear_filters = /** @type {((inputs?: Canvas_Clear_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Clear_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_clear_filters(inputs)
	return __es.canvas_clear_filters(inputs)
});
/**
* | output |
* | --- |
* | "Color" |
*
* @param {Canvas_Color_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_color_label = /** @type {((inputs?: Canvas_Color_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Color_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_color_label(inputs)
	return __es.canvas_color_label(inputs)
});
/**
* | output |
* | --- |
* | "Color mapping" |
*
* @param {Canvas_Color_MappingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_color_mapping = /** @type {((inputs?: Canvas_Color_MappingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Color_MappingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_color_mapping(inputs)
	return __es.canvas_color_mapping(inputs)
});
/**
* | output |
* | --- |
* | "Column dimensions" |
*
* @param {Canvas_Column_Dimensions_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_column_dimensions_label = /** @type {((inputs?: Canvas_Column_Dimensions_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Column_Dimensions_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_column_dimensions_label(inputs)
	return __es.canvas_column_dimensions_label(inputs)
});
/**
* | output |
* | --- |
* | "Columns" |
*
* @param {Canvas_Columns_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_columns_label = /** @type {((inputs?: Canvas_Columns_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Columns_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_columns_label(inputs)
	return __es.canvas_columns_label(inputs)
});
/**
* | output |
* | --- |
* | "Comparison values" |
*
* @param {Canvas_Comparison_Values_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_comparison_values_label = /** @type {((inputs?: Canvas_Comparison_Values_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Comparison_Values_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_comparison_values_label(inputs)
	return __es.canvas_comparison_values_label(inputs)
});
/**
* | output |
* | --- |
* | "{label} Configuration" |
*
* @param {Canvas_ConfigurationInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_configuration = /** @type {((inputs: Canvas_ConfigurationInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_ConfigurationInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_configuration(inputs)
	return __es.canvas_configuration(inputs)
});
/**
* | output |
* | --- |
* | "Canvas configurations" |
*
* @param {Canvas_ConfigurationsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_configurations = /** @type {((inputs?: Canvas_ConfigurationsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_ConfigurationsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_configurations(inputs)
	return __es.canvas_configurations(inputs)
});
/**
* | output |
* | --- |
* | "Custom Chart" |
*
* @param {Canvas_Custom_ChartInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_custom_chart = /** @type {((inputs?: Canvas_Custom_ChartInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Custom_ChartInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_custom_chart(inputs)
	return __es.canvas_custom_chart(inputs)
});
/**
* | output |
* | --- |
* | "Data labels" |
*
* @param {Canvas_Data_Labels_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_data_labels_label = /** @type {((inputs?: Canvas_Data_Labels_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Data_Labels_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_data_labels_label(inputs)
	return __es.canvas_data_labels_label(inputs)
});
/**
* | output |
* | --- |
* | "Delete" |
*
* @param {Canvas_DeleteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_delete = /** @type {((inputs?: Canvas_DeleteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_DeleteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_delete(inputs)
	return __es.canvas_delete(inputs)
});
/**
* | output |
* | --- |
* | "Delete widget?" |
*
* @param {Canvas_Delete_WidgetInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_delete_widget = /** @type {((inputs?: Canvas_Delete_WidgetInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Delete_WidgetInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_delete_widget(inputs)
	return __es.canvas_delete_widget(inputs)
});
/**
* | output |
* | --- |
* | "This widget and its configuration will be permanently removed." |
*
* @param {Canvas_Delete_Widget_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_delete_widget_description = /** @type {((inputs?: Canvas_Delete_Widget_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Delete_Widget_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_delete_widget_description(inputs)
	return __es.canvas_delete_widget_description(inputs)
});
/**
* | output |
* | --- |
* | "Describe chart changes..." |
*
* @param {Canvas_Describe_Chart_ChangesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_describe_chart_changes = /** @type {((inputs?: Canvas_Describe_Chart_ChangesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Describe_Chart_ChangesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_describe_chart_changes(inputs)
	return __es.canvas_describe_chart_changes(inputs)
});
/**
* | output |
* | --- |
* | "Describe the chart you want to see..." |
*
* @param {Canvas_Describe_Chart_PromptInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_describe_chart_prompt = /** @type {((inputs?: Canvas_Describe_Chart_PromptInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Describe_Chart_PromptInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_describe_chart_prompt(inputs)
	return __es.canvas_describe_chart_prompt(inputs)
});
/**
* | output |
* | --- |
* | "Description" |
*
* @param {Canvas_Description_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_description_label = /** @type {((inputs?: Canvas_Description_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Description_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_description_label(inputs)
	return __es.canvas_description_label(inputs)
});
/**
* | output |
* | --- |
* | "Add additional context for this component" |
*
* @param {Canvas_Description_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_description_placeholder = /** @type {((inputs?: Canvas_Description_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Description_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_description_placeholder(inputs)
	return __es.canvas_description_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Dimension" |
*
* @param {Canvas_Dimension_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_dimension_label = /** @type {((inputs?: Canvas_Dimension_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Dimension_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_dimension_label(inputs)
	return __es.canvas_dimension_label(inputs)
});
/**
* | output |
* | --- |
* | "Dimensions" |
*
* @param {Canvas_Dimensions_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_dimensions_label = /** @type {((inputs?: Canvas_Dimensions_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Dimensions_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_dimensions_label(inputs)
	return __es.canvas_dimensions_label(inputs)
});
/**
* | output |
* | --- |
* | "Display name" |
*
* @param {Canvas_Display_NameInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_display_name = /** @type {((inputs?: Canvas_Display_NameInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Display_NameInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_display_name(inputs)
	return __es.canvas_display_name(inputs)
});
/**
* | output |
* | --- |
* | "Diverging (Theme)" |
*
* @param {Canvas_Diverging_ThemeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_diverging_theme = /** @type {((inputs?: Canvas_Diverging_ThemeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Diverging_ThemeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_diverging_theme(inputs)
	return __es.canvas_diverging_theme(inputs)
});
/**
* | output |
* | --- |
* | "Edit" |
*
* @param {Canvas_EditInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_edit = /** @type {((inputs?: Canvas_EditInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_EditInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_edit(inputs)
	return __es.canvas_edit(inputs)
});
/**
* | output |
* | --- |
* | "Edit with AI" |
*
* @param {Canvas_Edit_With_AiInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_edit_with_ai = /** @type {((inputs?: Canvas_Edit_With_AiInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Edit_With_AiInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_edit_with_ai(inputs)
	return __es.canvas_edit_with_ai(inputs)
});
/**
* | output |
* | --- |
* | "End color" |
*
* @param {Canvas_End_ColorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_end_color = /** @type {((inputs?: Canvas_End_ColorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_End_ColorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_end_color(inputs)
	return __es.canvas_end_color(inputs)
});
/**
* | output |
* | --- |
* | "Enter a number" |
*
* @param {Canvas_Enter_A_NumberInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_enter_a_number = /** @type {((inputs?: Canvas_Enter_A_NumberInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Enter_A_NumberInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_enter_a_number(inputs)
	return __es.canvas_enter_a_number(inputs)
});
/**
* | output |
* | --- |
* | "Evenly distribute widgets" |
*
* @param {Canvas_Evenly_DistributeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_evenly_distribute = /** @type {((inputs?: Canvas_Evenly_DistributeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Evenly_DistributeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_evenly_distribute(inputs)
	return __es.canvas_evenly_distribute(inputs)
});
/**
* | output |
* | --- |
* | "Filter bar" |
*
* @param {Canvas_Filter_BarInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_filter_bar = /** @type {((inputs?: Canvas_Filter_BarInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Filter_BarInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_filter_bar(inputs)
	return __es.canvas_filter_bar(inputs)
});
/**
* | output |
* | --- |
* | "Filters" |
*
* @param {Canvas_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_filters = /** @type {((inputs?: Canvas_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_filters(inputs)
	return __es.canvas_filters(inputs)
});
/**
* | output |
* | --- |
* | "Format" |
*
* @param {Canvas_FormatInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_format = /** @type {((inputs?: Canvas_FormatInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_FormatInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_format(inputs)
	return __es.canvas_format(inputs)
});
/**
* | output |
* | --- |
* | "Generate" |
*
* @param {Canvas_GenerateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_generate = /** @type {((inputs?: Canvas_GenerateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_GenerateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_generate(inputs)
	return __es.canvas_generate(inputs)
});
/**
* | output |
* | --- |
* | "Gradient" |
*
* @param {Canvas_GradientInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_gradient = /** @type {((inputs?: Canvas_GradientInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_GradientInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_gradient(inputs)
	return __es.canvas_gradient(inputs)
});
/**
* | output |
* | --- |
* | "Greens" |
*
* @param {Canvas_GreensInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_greens = /** @type {((inputs?: Canvas_GreensInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_GreensInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_greens(inputs)
	return __es.canvas_greens(inputs)
});
/**
* | output |
* | --- |
* | "Greys" |
*
* @param {Canvas_GreysInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_greys = /** @type {((inputs?: Canvas_GreysInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_GreysInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_greys(inputs)
	return __es.canvas_greys(inputs)
});
/**
* | output |
* | --- |
* | "Hide labels below (%)" |
*
* @param {Canvas_Hide_Labels_BelowInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_hide_labels_below = /** @type {((inputs?: Canvas_Hide_Labels_BelowInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Hide_Labels_BelowInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_hide_labels_below(inputs)
	return __es.canvas_hide_labels_below(inputs)
});
/**
* | output |
* | --- |
* | "Hide total column" |
*
* @param {Canvas_Hide_Total_Column_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_hide_total_column_label = /** @type {((inputs?: Canvas_Hide_Total_Column_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Hide_Total_Column_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_hide_total_column_label(inputs)
	return __es.canvas_hide_total_column_label(inputs)
});
/**
* | output |
* | --- |
* | "Hide total row" |
*
* @param {Canvas_Hide_Total_Row_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_hide_total_row_label = /** @type {((inputs?: Canvas_Hide_Total_Row_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Hide_Total_Row_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_hide_total_row_label(inputs)
	return __es.canvas_hide_total_row_label(inputs)
});
/**
* | output |
* | --- |
* | "Image" |
*
* @param {Canvas_ImageInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_image = /** @type {((inputs?: Canvas_ImageInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_ImageInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_image(inputs)
	return __es.canvas_image(inputs)
});
/**
* | output |
* | --- |
* | "Inferno" |
*
* @param {Canvas_InfernoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_inferno = /** @type {((inputs?: Canvas_InfernoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_InfernoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_inferno(inputs)
	return __es.canvas_inferno(inputs)
});
/**
* | output |
* | --- |
* | "Inner Radius (%)" |
*
* @param {Canvas_Inner_Radius_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_inner_radius_label = /** @type {((inputs?: Canvas_Inner_Radius_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Inner_Radius_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_inner_radius_label(inputs)
	return __es.canvas_inner_radius_label(inputs)
});
/**
* | output |
* | --- |
* | "Insert widget" |
*
* @param {Canvas_Insert_WidgetInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_insert_widget = /** @type {((inputs?: Canvas_Insert_WidgetInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Insert_WidgetInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_insert_widget(inputs)
	return __es.canvas_insert_widget(inputs)
});
/**
* | output |
* | --- |
* | "Insert widget in row {row} at column {col}" |
*
* @param {Canvas_Insert_Widget_AtInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_insert_widget_at = /** @type {((inputs: Canvas_Insert_Widget_AtInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Insert_Widget_AtInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_insert_widget_at(inputs)
	return __es.canvas_insert_widget_at(inputs)
});
/**
* | output |
* | --- |
* | "Insert widget at column {col}" |
*
* @param {Canvas_Insert_Widget_At_ColInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_insert_widget_at_col = /** @type {((inputs: Canvas_Insert_Widget_At_ColInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Insert_Widget_At_ColInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_insert_widget_at_col(inputs)
	return __es.canvas_insert_widget_at_col(inputs)
});
/**
* | output |
* | --- |
* | "Insert widget in row {row}" |
*
* @param {Canvas_Insert_Widget_At_RowInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_insert_widget_at_row = /** @type {((inputs: Canvas_Insert_Widget_At_RowInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Insert_Widget_At_RowInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_insert_widget_at_row(inputs)
	return __es.canvas_insert_widget_at_row(inputs)
});
/**
* | output |
* | --- |
* | "KPI" |
*
* @param {Canvas_KpiInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_kpi = /** @type {((inputs?: Canvas_KpiInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_KpiInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_kpi(inputs)
	return __es.canvas_kpi(inputs)
});
/**
* | output |
* | --- |
* | "Label angle" |
*
* @param {Canvas_Label_AngleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_label_angle = /** @type {((inputs?: Canvas_Label_AngleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Label_AngleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_label_angle(inputs)
	return __es.canvas_label_angle(inputs)
});
/**
* | output |
* | --- |
* | "Leaderboard" |
*
* @param {Canvas_LeaderboardInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_leaderboard = /** @type {((inputs?: Canvas_LeaderboardInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_LeaderboardInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_leaderboard(inputs)
	return __es.canvas_leaderboard(inputs)
});
/**
* | output |
* | --- |
* | "Left Y-Axis" |
*
* @param {Canvas_Left_Y_Axis_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_left_y_axis_label = /** @type {((inputs?: Canvas_Left_Y_Axis_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Left_Y_Axis_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_left_y_axis_label(inputs)
	return __es.canvas_left_y_axis_label(inputs)
});
/**
* | output |
* | --- |
* | "Bottom" |
*
* @param {Canvas_Legend_BottomInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_legend_bottom = /** @type {((inputs?: Canvas_Legend_BottomInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Legend_BottomInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_legend_bottom(inputs)
	return __es.canvas_legend_bottom(inputs)
});
/**
* | output |
* | --- |
* | "Left" |
*
* @param {Canvas_Legend_LeftInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_legend_left = /** @type {((inputs?: Canvas_Legend_LeftInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Legend_LeftInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_legend_left(inputs)
	return __es.canvas_legend_left(inputs)
});
/**
* | output |
* | --- |
* | "None" |
*
* @param {Canvas_Legend_NoneInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_legend_none = /** @type {((inputs?: Canvas_Legend_NoneInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Legend_NoneInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_legend_none(inputs)
	return __es.canvas_legend_none(inputs)
});
/**
* | output |
* | --- |
* | "Legend orientation" |
*
* @param {Canvas_Legend_OrientationInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_legend_orientation = /** @type {((inputs?: Canvas_Legend_OrientationInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Legend_OrientationInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_legend_orientation(inputs)
	return __es.canvas_legend_orientation(inputs)
});
/**
* | output |
* | --- |
* | "Right" |
*
* @param {Canvas_Legend_RightInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_legend_right = /** @type {((inputs?: Canvas_Legend_RightInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Legend_RightInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_legend_right(inputs)
	return __es.canvas_legend_right(inputs)
});
/**
* | output |
* | --- |
* | "Top" |
*
* @param {Canvas_Legend_TopInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_legend_top = /** @type {((inputs?: Canvas_Legend_TopInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Legend_TopInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_legend_top(inputs)
	return __es.canvas_legend_top(inputs)
});
/**
* | output |
* | --- |
* | "Limit" |
*
* @param {Canvas_LimitInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_limit = /** @type {((inputs?: Canvas_LimitInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_LimitInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_limit(inputs)
	return __es.canvas_limit(inputs)
});
/**
* | output |
* | --- |
* | "Local filters" |
*
* @param {Canvas_Local_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_local_filters = /** @type {((inputs?: Canvas_Local_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Local_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_local_filters(inputs)
	return __es.canvas_local_filters(inputs)
});
/**
* | output |
* | --- |
* | "Local time range" |
*
* @param {Canvas_Local_Time_RangeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_local_time_range = /** @type {((inputs?: Canvas_Local_Time_RangeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Local_Time_RangeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_local_time_range(inputs)
	return __es.canvas_local_time_range(inputs)
});
/**
* | output |
* | --- |
* | "Magma" |
*
* @param {Canvas_MagmaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_magma = /** @type {((inputs?: Canvas_MagmaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_MagmaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_magma(inputs)
	return __es.canvas_magma(inputs)
});
/**
* | output |
* | --- |
* | "Mark type" |
*
* @param {Canvas_Mark_TypeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_mark_type = /** @type {((inputs?: Canvas_Mark_TypeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Mark_TypeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_mark_type(inputs)
	return __es.canvas_mark_type(inputs)
});
/**
* | output |
* | --- |
* | "Markdown" |
*
* @param {Canvas_Markdown_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_markdown_label = /** @type {((inputs?: Canvas_Markdown_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Markdown_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_markdown_label(inputs)
	return __es.canvas_markdown_label(inputs)
});
/**
* | output |
* | --- |
* | "Max" |
*
* @param {Canvas_MaxInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_max = /** @type {((inputs?: Canvas_MaxInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_MaxInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_max(inputs)
	return __es.canvas_max(inputs)
});
/**
* | output |
* | --- |
* | "Max width" |
*
* @param {Canvas_Max_WidthInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_max_width = /** @type {((inputs?: Canvas_Max_WidthInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Max_WidthInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_max_width(inputs)
	return __es.canvas_max_width(inputs)
});
/**
* | output |
* | --- |
* | "Measure" |
*
* @param {Canvas_Measure_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_measure_label = /** @type {((inputs?: Canvas_Measure_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Measure_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_measure_label(inputs)
	return __es.canvas_measure_label(inputs)
});
/**
* | output |
* | --- |
* | "Measures" |
*
* @param {Canvas_Measures_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_measures_label = /** @type {((inputs?: Canvas_Measures_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Measures_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_measures_label(inputs)
	return __es.canvas_measures_label(inputs)
});
/**
* | output |
* | --- |
* | "Metrics SQL" |
*
* @param {Canvas_Metrics_Sql_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_metrics_sql_label = /** @type {((inputs?: Canvas_Metrics_Sql_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Metrics_Sql_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_metrics_sql_label(inputs)
	return __es.canvas_metrics_sql_label(inputs)
});
/**
* | output |
* | --- |
* | "This will determine the measures and dimensions you can explore on this dashboard." |
*
* @param {Canvas_Metrics_View_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_metrics_view_description = /** @type {((inputs?: Canvas_Metrics_View_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Metrics_View_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_metrics_view_description(inputs)
	return __es.canvas_metrics_view_description(inputs)
});
/**
* | output |
* | --- |
* | "Metrics view" |
*
* @param {Canvas_Metrics_View_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_metrics_view_label = /** @type {((inputs?: Canvas_Metrics_View_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Metrics_View_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_metrics_view_label(inputs)
	return __es.canvas_metrics_view_label(inputs)
});
/**
* | output |
* | --- |
* | "Min" |
*
* @param {Canvas_MinInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_min = /** @type {((inputs?: Canvas_MinInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_MinInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_min(inputs)
	return __es.canvas_min(inputs)
});
/**
* | output |
* | --- |
* | "Mode" |
*
* @param {Canvas_Mode_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_mode_label = /** @type {((inputs?: Canvas_Mode_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Mode_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_mode_label(inputs)
	return __es.canvas_mode_label(inputs)
});
/**
* | output |
* | --- |
* | "Name" |
*
* @param {Canvas_Name_OptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_name_option = /** @type {((inputs?: Canvas_Name_OptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Name_OptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_name_option(inputs)
	return __es.canvas_name_option(inputs)
});
/**
* | output |
* | --- |
* | "No color values found" |
*
* @param {Canvas_No_Color_Values_FoundInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_no_color_values_found = /** @type {((inputs?: Canvas_No_Color_Values_FoundInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_No_Color_Values_FoundInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_no_color_values_found(inputs)
	return __es.canvas_no_color_values_found(inputs)
});
/**
* | output |
* | --- |
* | "No components added" |
*
* @param {Canvas_No_Components_AddedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_no_components_added = /** @type {((inputs?: Canvas_No_Components_AddedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_No_Components_AddedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_no_components_added(inputs)
	return __es.canvas_no_components_added(inputs)
});
/**
* | output |
* | --- |
* | "No filters selected" |
*
* @param {Canvas_No_Filters_SelectedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_no_filters_selected = /** @type {((inputs?: Canvas_No_Filters_SelectedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_No_Filters_SelectedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_no_filters_selected(inputs)
	return __es.canvas_no_filters_selected(inputs)
});
/**
* | output |
* | --- |
* | "No valid component {id} in project" |
*
* @param {Canvas_No_Valid_ComponentInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_no_valid_component = /** @type {((inputs: Canvas_No_Valid_ComponentInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_No_Valid_ComponentInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_no_valid_component(inputs)
	return __es.canvas_no_valid_component(inputs)
});
/**
* | output |
* | --- |
* | "No valid metrics view in project" |
*
* @param {Canvas_No_Valid_Metrics_ViewInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_no_valid_metrics_view = /** @type {((inputs?: Canvas_No_Valid_Metrics_ViewInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_No_Valid_Metrics_ViewInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_no_valid_metrics_view(inputs)
	return __es.canvas_no_valid_metrics_view(inputs)
});
/**
* | output |
* | --- |
* | "Canvas not found" |
*
* @param {Canvas_Not_FoundInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_not_found = /** @type {((inputs?: Canvas_Not_FoundInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Not_FoundInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_not_found(inputs)
	return __es.canvas_not_found(inputs)
});
/**
* | output |
* | --- |
* | "Number of rows" |
*
* @param {Canvas_Number_Of_Rows_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_number_of_rows_label = /** @type {((inputs?: Canvas_Number_Of_Rows_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Number_Of_Rows_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_number_of_rows_label(inputs)
	return __es.canvas_number_of_rows_label(inputs)
});
/**
* | output |
* | --- |
* | "Oranges" |
*
* @param {Canvas_OrangesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_oranges = /** @type {((inputs?: Canvas_OrangesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_OrangesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_oranges(inputs)
	return __es.canvas_oranges(inputs)
});
/**
* | output |
* | --- |
* | "Order" |
*
* @param {Canvas_Order_OptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_order_option = /** @type {((inputs?: Canvas_Order_OptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Order_OptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_order_option(inputs)
	return __es.canvas_order_option(inputs)
});
/**
* | output |
* | --- |
* | "Percent of" |
*
* @param {Canvas_Percent_Of_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_percent_of_label = /** @type {((inputs?: Canvas_Percent_Of_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Percent_Of_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_percent_of_label(inputs)
	return __es.canvas_percent_of_label(inputs)
});
/**
* | output |
* | --- |
* | "Plasma" |
*
* @param {Canvas_PlasmaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_plasma = /** @type {((inputs?: Canvas_PlasmaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_PlasmaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_plasma(inputs)
	return __es.canvas_plasma(inputs)
});
/**
* | output |
* | --- |
* | "Previous" |
*
* @param {Canvas_Previous_OptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_previous_option = /** @type {((inputs?: Canvas_Previous_OptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Previous_OptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_previous_option(inputs)
	return __es.canvas_previous_option(inputs)
});
/**
* | output |
* | --- |
* | "Purples" |
*
* @param {Canvas_PurplesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_purples = /** @type {((inputs?: Canvas_PurplesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_PurplesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_purples(inputs)
	return __es.canvas_purples(inputs)
});
/**
* | output |
* | --- |
* | "Reds" |
*
* @param {Canvas_RedsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_reds = /** @type {((inputs?: Canvas_RedsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_RedsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_reds(inputs)
	return __es.canvas_reds(inputs)
});
/**
* | output |
* | --- |
* | "Remove query {idx}" |
*
* @param {Canvas_Remove_Query_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_remove_query_aria = /** @type {((inputs: Canvas_Remove_Query_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Remove_Query_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_remove_query_aria(inputs)
	return __es.canvas_remove_query_aria(inputs)
});
/**
* | output |
* | --- |
* | "Reset to default" |
*
* @param {Canvas_Reset_To_DefaultInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_reset_to_default = /** @type {((inputs?: Canvas_Reset_To_DefaultInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Reset_To_DefaultInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_reset_to_default(inputs)
	return __es.canvas_reset_to_default(inputs)
});
/**
* | output |
* | --- |
* | "Resize row {row} column {col}" |
*
* @param {Canvas_Resize_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_resize_aria = /** @type {((inputs: Canvas_Resize_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Resize_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_resize_aria(inputs)
	return __es.canvas_resize_aria(inputs)
});
/**
* | output |
* | --- |
* | "Right Y-Axis" |
*
* @param {Canvas_Right_Y_Axis_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_right_y_axis_label = /** @type {((inputs?: Canvas_Right_Y_Axis_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Right_Y_Axis_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_right_y_axis_label(inputs)
	return __es.canvas_right_y_axis_label(inputs)
});
/**
* | output |
* | --- |
* | "Row dimensions" |
*
* @param {Canvas_Row_Dimensions_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_row_dimensions_label = /** @type {((inputs?: Canvas_Row_Dimensions_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Row_Dimensions_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_row_dimensions_label(inputs)
	return __es.canvas_row_dimensions_label(inputs)
});
/**
* | output |
* | --- |
* | "Save as default" |
*
* @param {Canvas_Save_As_DefaultInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_save_as_default = /** @type {((inputs?: Canvas_Save_As_DefaultInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Save_As_DefaultInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_save_as_default(inputs)
	return __es.canvas_save_as_default(inputs)
});
/**
* | output |
* | --- |
* | "Saved default filters" |
*
* @param {Canvas_Saved_Default_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_saved_default_filters = /** @type {((inputs?: Canvas_Saved_Default_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Saved_Default_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_saved_default_filters(inputs)
	return __es.canvas_saved_default_filters(inputs)
});
/**
* | output |
* | --- |
* | "Saving default filters" |
*
* @param {Canvas_Saving_Default_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_saving_default_filters = /** @type {((inputs?: Canvas_Saving_Default_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Saving_Default_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_saving_default_filters(inputs)
	return __es.canvas_saving_default_filters(inputs)
});
/**
* | output |
* | --- |
* | "Scheme" |
*
* @param {Canvas_SchemeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_scheme = /** @type {((inputs?: Canvas_SchemeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_SchemeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_scheme(inputs)
	return __es.canvas_scheme(inputs)
});
/**
* | output |
* | --- |
* | "See less" |
*
* @param {Canvas_See_LessInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_see_less = /** @type {((inputs?: Canvas_See_LessInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_See_LessInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_see_less(inputs)
	return __es.canvas_see_less(inputs)
});
/**
* | output |
* | --- |
* | "See 1 more value" |
*
* @param {Canvas_See_More_ValueInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_see_more_value = /** @type {((inputs?: Canvas_See_More_ValueInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_See_More_ValueInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_see_more_value(inputs)
	return __es.canvas_see_more_value(inputs)
});
/**
* | output |
* | --- |
* | "See {count} more values" |
*
* @param {Canvas_See_More_ValuesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_see_more_values = /** @type {((inputs: Canvas_See_More_ValuesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_See_More_ValuesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_see_more_values(inputs)
	return __es.canvas_see_more_values(inputs)
});
/**
* | output |
* | --- |
* | "Select a metrics view" |
*
* @param {Canvas_Select_Metrics_ViewInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_select_metrics_view = /** @type {((inputs?: Canvas_Select_Metrics_ViewInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Select_Metrics_ViewInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_select_metrics_view(inputs)
	return __es.canvas_select_metrics_view(inputs)
});
/**
* | output |
* | --- |
* | "Sequential (Theme)" |
*
* @param {Canvas_Sequential_ThemeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_sequential_theme = /** @type {((inputs?: Canvas_Sequential_ThemeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Sequential_ThemeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_sequential_theme(inputs)
	return __es.canvas_sequential_theme(inputs)
});
/**
* | output |
* | --- |
* | "Show axis title" |
*
* @param {Canvas_Show_Axis_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_show_axis_title = /** @type {((inputs?: Canvas_Show_Axis_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Show_Axis_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_show_axis_title(inputs)
	return __es.canvas_show_axis_title(inputs)
});
/**
* | output |
* | --- |
* | "Show description as tooltip" |
*
* @param {Canvas_Show_Description_As_Tooltip_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_show_description_as_tooltip_label = /** @type {((inputs?: Canvas_Show_Description_As_Tooltip_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Show_Description_As_Tooltip_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_show_description_as_tooltip_label(inputs)
	return __es.canvas_show_description_as_tooltip_label(inputs)
});
/**
* | output |
* | --- |
* | "Show null values" |
*
* @param {Canvas_Show_Null_ValuesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_show_null_values = /** @type {((inputs?: Canvas_Show_Null_ValuesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Show_Null_ValuesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_show_null_values(inputs)
	return __es.canvas_show_null_values(inputs)
});
/**
* | output |
* | --- |
* | "Show \"Other\" bucket" |
*
* @param {Canvas_Show_Other_Bucket_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_show_other_bucket_label = /** @type {((inputs?: Canvas_Show_Other_Bucket_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Show_Other_Bucket_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_show_other_bucket_label(inputs)
	return __es.canvas_show_other_bucket_label(inputs)
});
/**
* | output |
* | --- |
* | "Show totals value" |
*
* @param {Canvas_Show_Totals_ValueInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_show_totals_value = /** @type {((inputs?: Canvas_Show_Totals_ValueInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Show_Totals_ValueInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_show_totals_value(inputs)
	return __es.canvas_show_totals_value(inputs)
});
/**
* | output |
* | --- |
* | "Size" |
*
* @param {Canvas_Size_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_size_label = /** @type {((inputs?: Canvas_Size_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Size_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_size_label(inputs)
	return __es.canvas_size_label(inputs)
});
/**
* | output |
* | --- |
* | "Sort" |
*
* @param {Canvas_SortInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_sort = /** @type {((inputs?: Canvas_SortInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_SortInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_sort(inputs)
	return __es.canvas_sort(inputs)
});
/**
* | output |
* | --- |
* | "Show sparkline below the value" |
*
* @param {Canvas_Sparkline_BelowInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_sparkline_below = /** @type {((inputs?: Canvas_Sparkline_BelowInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Sparkline_BelowInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_sparkline_below(inputs)
	return __es.canvas_sparkline_below(inputs)
});
/**
* | output |
* | --- |
* | "Sparkline" |
*
* @param {Canvas_Sparkline_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_sparkline_label = /** @type {((inputs?: Canvas_Sparkline_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Sparkline_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_sparkline_label(inputs)
	return __es.canvas_sparkline_label(inputs)
});
/**
* | output |
* | --- |
* | "Show sparkline to the right of the value" |
*
* @param {Canvas_Sparkline_RightInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_sparkline_right = /** @type {((inputs?: Canvas_Sparkline_RightInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Sparkline_RightInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_sparkline_right(inputs)
	return __es.canvas_sparkline_right(inputs)
});
/**
* | output |
* | --- |
* | "Spectral" |
*
* @param {Canvas_SpectralInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_spectral = /** @type {((inputs?: Canvas_SpectralInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_SpectralInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_spectral(inputs)
	return __es.canvas_spectral(inputs)
});
/**
* | output |
* | --- |
* | "Stage" |
*
* @param {Canvas_Stage_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_stage_label = /** @type {((inputs?: Canvas_Stage_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Stage_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_stage_label(inputs)
	return __es.canvas_stage_label(inputs)
});
/**
* | output |
* | --- |
* | "Start color" |
*
* @param {Canvas_Start_ColorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_start_color = /** @type {((inputs?: Canvas_Start_ColorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Start_ColorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_start_color(inputs)
	return __es.canvas_start_color(inputs)
});
/**
* | output |
* | --- |
* | "Switch to {mark} editor" |
*
* @param {Canvas_Switch_Mark_Type_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_switch_mark_type_aria = /** @type {((inputs: Canvas_Switch_Mark_Type_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Switch_Mark_Type_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_switch_mark_type_aria(inputs)
	return __es.canvas_switch_mark_type_aria(inputs)
});
/**
* | output |
* | --- |
* | "Table" |
*
* @param {Canvas_TableInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_table = /** @type {((inputs?: Canvas_TableInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_TableInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_table(inputs)
	return __es.canvas_table(inputs)
});
/**
* | output |
* | --- |
* | "Table type" |
*
* @param {Canvas_Table_TypeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_table_type = /** @type {((inputs?: Canvas_Table_TypeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Table_TypeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_table_type(inputs)
	return __es.canvas_table_type(inputs)
});
/**
* | output |
* | --- |
* | "Teal blues" |
*
* @param {Canvas_Teal_BluesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_teal_blues = /** @type {((inputs?: Canvas_Teal_BluesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Teal_BluesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_teal_blues(inputs)
	return __es.canvas_teal_blues(inputs)
});
/**
* | output |
* | --- |
* | "Teals" |
*
* @param {Canvas_TealsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_teals = /** @type {((inputs?: Canvas_TealsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_TealsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_teals(inputs)
	return __es.canvas_teals(inputs)
});
/**
* | output |
* | --- |
* | "Text/Markdown" |
*
* @param {Canvas_Text_MarkdownInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_text_markdown = /** @type {((inputs?: Canvas_Text_MarkdownInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Text_MarkdownInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_text_markdown(inputs)
	return __es.canvas_text_markdown(inputs)
});
/**
* | output |
* | --- |
* | "Time range display" |
*
* @param {Canvas_Time_Range_Display_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_time_range_display_label = /** @type {((inputs?: Canvas_Time_Range_Display_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Time_Range_Display_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_time_range_display_label(inputs)
	return __es.canvas_time_range_display_label(inputs)
});
/**
* | output |
* | --- |
* | "Time ranges" |
*
* @param {Canvas_Time_RangesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_time_ranges = /** @type {((inputs?: Canvas_Time_RangesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Time_RangesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_time_ranges(inputs)
	return __es.canvas_time_ranges(inputs)
});
/**
* | output |
* | --- |
* | "Time zones" |
*
* @param {Canvas_Time_ZonesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_time_zones = /** @type {((inputs?: Canvas_Time_ZonesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Time_ZonesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_time_zones(inputs)
	return __es.canvas_time_zones(inputs)
});
/**
* | output |
* | --- |
* | "Title" |
*
* @param {Canvas_Title_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_title_label = /** @type {((inputs?: Canvas_Title_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Title_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_title_label(inputs)
	return __es.canvas_title_label(inputs)
});
/**
* | output |
* | --- |
* | "Add a title to describe this component" |
*
* @param {Canvas_Title_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_title_placeholder = /** @type {((inputs?: Canvas_Title_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Title_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_title_placeholder(inputs)
	return __es.canvas_title_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Tooltip" |
*
* @param {Canvas_Tooltip_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_tooltip_label = /** @type {((inputs?: Canvas_Tooltip_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Tooltip_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_tooltip_label(inputs)
	return __es.canvas_tooltip_label(inputs)
});
/**
* | output |
* | --- |
* | "Top" |
*
* @param {Canvas_Top_OptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_top_option = /** @type {((inputs?: Canvas_Top_OptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Top_OptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_top_option(inputs)
	return __es.canvas_top_option(inputs)
});
/**
* | output |
* | --- |
* | "Turbo" |
*
* @param {Canvas_TurboInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_turbo = /** @type {((inputs?: Canvas_TurboInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_TurboInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_turbo(inputs)
	return __es.canvas_turbo(inputs)
});
/**
* | output |
* | --- |
* | "An unknown error occurred." |
*
* @param {Canvas_Unknown_ErrorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_unknown_error = /** @type {((inputs?: Canvas_Unknown_ErrorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Unknown_ErrorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_unknown_error(inputs)
	return __es.canvas_unknown_error(inputs)
});
/**
* | output |
* | --- |
* | "URL" |
*
* @param {Canvas_Url_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_url_label = /** @type {((inputs?: Canvas_Url_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Url_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_url_label(inputs)
	return __es.canvas_url_label(inputs)
});
/**
* | output |
* | --- |
* | "Value" |
*
* @param {Canvas_Value_OptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_value_option = /** @type {((inputs?: Canvas_Value_OptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Value_OptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_value_option(inputs)
	return __es.canvas_value_option(inputs)
});
/**
* | output |
* | --- |
* | "Vega Lite Spec" |
*
* @param {Canvas_Vega_Lite_Spec_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_vega_lite_spec_label = /** @type {((inputs?: Canvas_Vega_Lite_Spec_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Vega_Lite_Spec_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_vega_lite_spec_label(inputs)
	return __es.canvas_vega_lite_spec_label(inputs)
});
/**
* | output |
* | --- |
* | "Viewing default state" |
*
* @param {Canvas_Viewing_Default_StateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_viewing_default_state = /** @type {((inputs?: Canvas_Viewing_Default_StateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Viewing_Default_StateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_viewing_default_state(inputs)
	return __es.canvas_viewing_default_state(inputs)
});
/**
* | output |
* | --- |
* | "Viridis" |
*
* @param {Canvas_ViridisInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_viridis = /** @type {((inputs?: Canvas_ViridisInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_ViridisInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_viridis(inputs)
	return __es.canvas_viridis(inputs)
});
/**
* | output |
* | --- |
* | "Which metrics view should this dashboard reference?" |
*
* @param {Canvas_Which_Metrics_ViewInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_which_metrics_view = /** @type {((inputs?: Canvas_Which_Metrics_ViewInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Which_Metrics_ViewInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_which_metrics_view(inputs)
	return __es.canvas_which_metrics_view(inputs)
});
/**
* | output |
* | --- |
* | "Width" |
*
* @param {Canvas_Width_OptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_width_option = /** @type {((inputs?: Canvas_Width_OptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Width_OptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_width_option(inputs)
	return __es.canvas_width_option(inputs)
});
/**
* | output |
* | --- |
* | "X-axis" |
*
* @param {Canvas_X_Axis_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_x_axis_label = /** @type {((inputs?: Canvas_X_Axis_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_X_Axis_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_x_axis_label(inputs)
	return __es.canvas_x_axis_label(inputs)
});
/**
* | output |
* | --- |
* | "Y-axis" |
*
* @param {Canvas_Y_Axis_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_y_axis_label = /** @type {((inputs?: Canvas_Y_Axis_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Y_Axis_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_y_axis_label(inputs)
	return __es.canvas_y_axis_label(inputs)
});
/**
* | output |
* | --- |
* | "Zero based origin" |
*
* @param {Canvas_Zero_Based_OriginInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const canvas_zero_based_origin = /** @type {((inputs?: Canvas_Zero_Based_OriginInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Canvas_Zero_Based_OriginInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.canvas_zero_based_origin(inputs)
	return __es.canvas_zero_based_origin(inputs)
});
/**
* | output |
* | --- |
* | "Adaptive" |
*
* @param {Chart_AdaptiveInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_adaptive = /** @type {((inputs?: Chart_AdaptiveInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_AdaptiveInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_adaptive(inputs)
	return __es.chart_adaptive(inputs)
});
/**
* | output |
* | --- |
* | "Adaptive: Line chart by default. Switches to bar when there are few data points, and stacked bar when comparing dimension" |
*
* @param {Chart_Adaptive_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_adaptive_tooltip = /** @type {((inputs?: Chart_Adaptive_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_Adaptive_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_adaptive_tooltip(inputs)
	return __es.chart_adaptive_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "Add comparison values to use {chartType} chart" |
*
* @param {Chart_Add_Comparison_To_UseInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_add_comparison_to_use = /** @type {((inputs: Chart_Add_Comparison_To_UseInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_Add_Comparison_To_UseInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_add_comparison_to_use(inputs)
	return __es.chart_add_comparison_to_use(inputs)
});
/**
* | output |
* | --- |
* | "Bar" |
*
* @param {Chart_BarInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_bar = /** @type {((inputs?: Chart_BarInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_BarInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_bar(inputs)
	return __es.chart_bar(inputs)
});
/**
* | output |
* | --- |
* | "Copy" |
*
* @param {Chart_Copy_To_ClipboardInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_copy_to_clipboard = /** @type {((inputs?: Chart_Copy_To_ClipboardInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_Copy_To_ClipboardInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_copy_to_clipboard(inputs)
	return __es.chart_copy_to_clipboard(inputs)
});
/**
* | output |
* | --- |
* | "Line" |
*
* @param {Chart_LineInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_line = /** @type {((inputs?: Chart_LineInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_LineInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_line(inputs)
	return __es.chart_line(inputs)
});
/**
* | output |
* | --- |
* | "Stacked area" |
*
* @param {Chart_Stacked_AreaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_stacked_area = /** @type {((inputs?: Chart_Stacked_AreaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_Stacked_AreaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_stacked_area(inputs)
	return __es.chart_stacked_area(inputs)
});
/**
* | output |
* | --- |
* | "Stacked bar" |
*
* @param {Chart_Stacked_BarInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_stacked_bar = /** @type {((inputs?: Chart_Stacked_BarInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_Stacked_BarInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_stacked_bar(inputs)
	return __es.chart_stacked_bar(inputs)
});
/**
* | output |
* | --- |
* | "Undo zoom" |
*
* @param {Chart_Undo_ZoomInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_undo_zoom = /** @type {((inputs?: Chart_Undo_ZoomInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_Undo_ZoomInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_undo_zoom(inputs)
	return __es.chart_undo_zoom(inputs)
});
/**
* | output |
* | --- |
* | "Undo Zoom" |
*
* @param {Chart_Undo_Zoom_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_undo_zoom_label = /** @type {((inputs?: Chart_Undo_Zoom_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_Undo_Zoom_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_undo_zoom_label(inputs)
	return __es.chart_undo_zoom_label(inputs)
});
/**
* | output |
* | --- |
* | "vs" |
*
* @param {Chart_VsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_vs = /** @type {((inputs?: Chart_VsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_VsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_vs(inputs)
	return __es.chart_vs(inputs)
});
/**
* | output |
* | --- |
* | "Zoom" |
*
* @param {Chart_ZoomInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_zoom = /** @type {((inputs?: Chart_ZoomInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_ZoomInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_zoom(inputs)
	return __es.chart_zoom(inputs)
});
/**
* | output |
* | --- |
* | "Zoom" |
*
* @param {Chart_Zoom_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chart_zoom_label = /** @type {((inputs?: Chart_Zoom_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chart_Zoom_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chart_zoom_label(inputs)
	return __es.chart_zoom_label(inputs)
});
/**
* | output |
* | --- |
* | "AI can make mistakes. Consider your dashboard the source of truth." |
*
* @param {Chat_Ai_DisclaimerInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_ai_disclaimer = /** @type {((inputs?: Chat_Ai_DisclaimerInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Ai_DisclaimerInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_ai_disclaimer(inputs)
	return __es.chat_ai_disclaimer(inputs)
});
/**
* | output |
* | --- |
* | "Cancel streaming" |
*
* @param {Chat_Cancel_StreamingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_cancel_streaming = /** @type {((inputs?: Chat_Cancel_StreamingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Cancel_StreamingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_cancel_streaming(inputs)
	return __es.chat_cancel_streaming(inputs)
});
/**
* | output |
* | --- |
* | "Close chat" |
*
* @param {Chat_CloseInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_close = /** @type {((inputs?: Chat_CloseInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_CloseInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_close(inputs)
	return __es.chat_close(inputs)
});
/**
* | output |
* | --- |
* | "Connect your own client" |
*
* @param {Chat_Connect_ClientInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_connect_client = /** @type {((inputs?: Chat_Connect_ClientInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Connect_ClientInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_connect_client(inputs)
	return __es.chat_connect_client(inputs)
});
/**
* | output |
* | --- |
* | "Conversation history" |
*
* @param {Chat_Conversation_HistoryInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_conversation_history = /** @type {((inputs?: Chat_Conversation_HistoryInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Conversation_HistoryInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_conversation_history(inputs)
	return __es.chat_conversation_history(inputs)
});
/**
* | output |
* | --- |
* | "Downvote response" |
*
* @param {Chat_Downvote_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_downvote_aria = /** @type {((inputs?: Chat_Downvote_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Downvote_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_downvote_aria(inputs)
	return __es.chat_downvote_aria(inputs)
});
/**
* | output |
* | --- |
* | "This response needs improvement" |
*
* @param {Chat_Downvote_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_downvote_tooltip = /** @type {((inputs?: Chat_Downvote_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Downvote_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_downvote_tooltip(inputs)
	return __es.chat_downvote_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "less than a second" |
*
* @param {Chat_Duration_Less_Than_SecondInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_duration_less_than_second = /** @type {((inputs?: Chat_Duration_Less_Than_SecondInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Duration_Less_Than_SecondInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_duration_less_than_second(inputs)
	return __es.chat_duration_less_than_second(inputs)
});
/**
* | output |
* | --- |
* | "{count} minutes" |
*
* @param {Chat_Duration_MinutesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_duration_minutes = /** @type {((inputs: Chat_Duration_MinutesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Duration_MinutesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_duration_minutes(inputs)
	return __es.chat_duration_minutes(inputs)
});
/**
* | output |
* | --- |
* | "1 minute" |
*
* @param {Chat_Duration_One_MinuteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_duration_one_minute = /** @type {((inputs?: Chat_Duration_One_MinuteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Duration_One_MinuteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_duration_one_minute(inputs)
	return __es.chat_duration_one_minute(inputs)
});
/**
* | output |
* | --- |
* | "1 second" |
*
* @param {Chat_Duration_One_SecondInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_duration_one_second = /** @type {((inputs?: Chat_Duration_One_SecondInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Duration_One_SecondInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_duration_one_second(inputs)
	return __es.chat_duration_one_second(inputs)
});
/**
* | output |
* | --- |
* | "{count} seconds" |
*
* @param {Chat_Duration_SecondsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_duration_seconds = /** @type {((inputs: Chat_Duration_SecondsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Duration_SecondsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_duration_seconds(inputs)
	return __es.chat_duration_seconds(inputs)
});
/**
* | output |
* | --- |
* | "Happy to help explore your data" |
*
* @param {Chat_Empty_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_empty_label = /** @type {((inputs?: Chat_Empty_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Empty_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_empty_label(inputs)
	return __es.chat_empty_label(inputs)
});
/**
* | output |
* | --- |
* | "Failed to generate response" |
*
* @param {Chat_Failed_To_GenerateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_failed_to_generate = /** @type {((inputs?: Chat_Failed_To_GenerateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Failed_To_GenerateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_failed_to_generate(inputs)
	return __es.chat_failed_to_generate(inputs)
});
/**
* | output |
* | --- |
* | "Analyzing feedback..." |
*
* @param {Chat_Feedback_AnalyzingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_feedback_analyzing = /** @type {((inputs?: Chat_Feedback_AnalyzingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Feedback_AnalyzingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_feedback_analyzing(inputs)
	return __es.chat_feedback_analyzing(inputs)
});
/**
* | output |
* | --- |
* | "Comments" |
*
* @param {Chat_Feedback_CommentsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_feedback_comments = /** @type {((inputs?: Chat_Feedback_CommentsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Feedback_CommentsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_feedback_comments(inputs)
	return __es.chat_feedback_comments(inputs)
});
/**
* | output |
* | --- |
* | "Type here..." |
*
* @param {Chat_Feedback_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_feedback_placeholder = /** @type {((inputs?: Chat_Feedback_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Feedback_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_feedback_placeholder(inputs)
	return __es.chat_feedback_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Select all that apply." |
*
* @param {Chat_Feedback_Select_AllInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_feedback_select_all = /** @type {((inputs?: Chat_Feedback_Select_AllInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Feedback_Select_AllInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_feedback_select_all(inputs)
	return __es.chat_feedback_select_all(inputs)
});
/**
* | output |
* | --- |
* | "Skip" |
*
* @param {Chat_Feedback_SkipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_feedback_skip = /** @type {((inputs?: Chat_Feedback_SkipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Feedback_SkipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_feedback_skip(inputs)
	return __es.chat_feedback_skip(inputs)
});
/**
* | output |
* | --- |
* | "Submit" |
*
* @param {Chat_Feedback_SubmitInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_feedback_submit = /** @type {((inputs?: Chat_Feedback_SubmitInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Feedback_SubmitInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_feedback_submit(inputs)
	return __es.chat_feedback_submit(inputs)
});
/**
* | output |
* | --- |
* | "Give feedback" |
*
* @param {Chat_Feedback_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_feedback_title = /** @type {((inputs?: Chat_Feedback_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Feedback_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_feedback_title(inputs)
	return __es.chat_feedback_title(inputs)
});
/**
* | output |
* | --- |
* | "{days}d ago" |
*
* @param {Chat_Group_Days_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_group_days_ago = /** @type {((inputs: Chat_Group_Days_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Group_Days_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_group_days_ago(inputs)
	return __es.chat_group_days_ago(inputs)
});
/**
* | output |
* | --- |
* | "Older" |
*
* @param {Chat_Group_OlderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_group_older = /** @type {((inputs?: Chat_Group_OlderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Group_OlderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_group_older(inputs)
	return __es.chat_group_older(inputs)
});
/**
* | output |
* | --- |
* | "Today" |
*
* @param {Chat_Group_TodayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_group_today = /** @type {((inputs?: Chat_Group_TodayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Group_TodayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_group_today(inputs)
	return __es.chat_group_today(inputs)
});
/**
* | output |
* | --- |
* | "Yesterday" |
*
* @param {Chat_Group_YesterdayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_group_yesterday = /** @type {((inputs?: Chat_Group_YesterdayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Group_YesterdayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_group_yesterday(inputs)
	return __es.chat_group_yesterday(inputs)
});
/**
* | output |
* | --- |
* | "Happy to help explore your data" |
*
* @param {Chat_Happy_To_ExploreInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_happy_to_explore = /** @type {((inputs?: Chat_Happy_To_ExploreInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Happy_To_ExploreInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_happy_to_explore(inputs)
	return __es.chat_happy_to_explore(inputs)
});
/**
* | output |
* | --- |
* | "How can I help you today?" |
*
* @param {Chat_How_Can_I_HelpInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_how_can_i_help = /** @type {((inputs?: Chat_How_Can_I_HelpInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_How_Can_I_HelpInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_how_can_i_help(inputs)
	return __es.chat_how_can_i_help(inputs)
});
/**
* | output |
* | --- |
* | "New conversation" |
*
* @param {Chat_New_ConversationInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_new_conversation = /** @type {((inputs?: Chat_New_ConversationInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_New_ConversationInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_new_conversation(inputs)
	return __es.chat_new_conversation(inputs)
});
/**
* | output |
* | --- |
* | "No conversations yet." |
*
* @param {Chat_No_ConversationsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_no_conversations = /** @type {((inputs?: Chat_No_ConversationsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_No_ConversationsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_no_conversations(inputs)
	return __es.chat_no_conversations(inputs)
});
/**
* | output |
* | --- |
* | "Type a question, or press @ to insert a metric, dimension, or measure." |
*
* @param {Chat_Placeholder_AnalystInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_placeholder_analyst = /** @type {((inputs?: Chat_Placeholder_AnalystInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Placeholder_AnalystInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_placeholder_analyst(inputs)
	return __es.chat_placeholder_analyst(inputs)
});
/**
* | output |
* | --- |
* | "Send message" |
*
* @param {Chat_Send_MessageInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_send_message = /** @type {((inputs?: Chat_Send_MessageInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Send_MessageInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_send_message(inputs)
	return __es.chat_send_message(inputs)
});
/**
* | output |
* | --- |
* | "Share conversation" |
*
* @param {Chat_Share_ConversationInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_share_conversation = /** @type {((inputs?: Chat_Share_ConversationInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Share_ConversationInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_share_conversation(inputs)
	return __es.chat_share_conversation(inputs)
});
/**
* | output |
* | --- |
* | "Copied!" |
*
* @param {Chat_Share_CopiedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_share_copied = /** @type {((inputs?: Chat_Share_CopiedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Share_CopiedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_share_copied(inputs)
	return __es.chat_share_copied(inputs)
});
/**
* | output |
* | --- |
* | "Create link" |
*
* @param {Chat_Share_Create_LinkInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_share_create_link = /** @type {((inputs?: Chat_Share_Create_LinkInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Share_Create_LinkInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_share_create_link(inputs)
	return __es.chat_share_create_link(inputs)
});
/**
* | output |
* | --- |
* | "Creating link..." |
*
* @param {Chat_Share_CreatingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_share_creating = /** @type {((inputs?: Chat_Share_CreatingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Share_CreatingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_share_creating(inputs)
	return __es.chat_share_creating(inputs)
});
/**
* | output |
* | --- |
* | "Share this conversation with other project members. They can view and continue the conversation." |
*
* @param {Chat_Share_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_share_description = /** @type {((inputs?: Chat_Share_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Share_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_share_description(inputs)
	return __es.chat_share_description(inputs)
});
/**
* | output |
* | --- |
* | "Start a conversation to share" |
*
* @param {Chat_Share_Start_FirstInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_share_start_first = /** @type {((inputs?: Chat_Share_Start_FirstInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Share_Start_FirstInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_share_start_first(inputs)
	return __es.chat_share_start_first(inputs)
});
/**
* | output |
* | --- |
* | "Show details" |
*
* @param {Chat_Show_DetailsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_show_details = /** @type {((inputs?: Chat_Show_DetailsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Show_DetailsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_show_details(inputs)
	return __es.chat_show_details(inputs)
});
/**
* | output |
* | --- |
* | "Thinking" |
*
* @param {Chat_ThinkingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_thinking = /** @type {((inputs?: Chat_ThinkingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_ThinkingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_thinking(inputs)
	return __es.chat_thinking(inputs)
});
/**
* | output |
* | --- |
* | "Thought for {duration}" |
*
* @param {Chat_Thought_ForInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_thought_for = /** @type {((inputs: Chat_Thought_ForInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Thought_ForInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_thought_for(inputs)
	return __es.chat_thought_for(inputs)
});
/**
* | output |
* | --- |
* | "Unable to load conversation" |
*
* @param {Chat_Unable_To_LoadInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_unable_to_load = /** @type {((inputs?: Chat_Unable_To_LoadInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Unable_To_LoadInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_unable_to_load(inputs)
	return __es.chat_unable_to_load(inputs)
});
/**
* | output |
* | --- |
* | "Upvote response" |
*
* @param {Chat_Upvote_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_upvote_aria = /** @type {((inputs?: Chat_Upvote_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Upvote_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_upvote_aria(inputs)
	return __es.chat_upvote_aria(inputs)
});
/**
* | output |
* | --- |
* | "This response was helpful" |
*
* @param {Chat_Upvote_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const chat_upvote_tooltip = /** @type {((inputs?: Chat_Upvote_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Chat_Upvote_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.chat_upvote_tooltip(inputs)
	return __es.chat_upvote_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "Apply" |
*
* @param {Common_ApplyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const common_apply = /** @type {((inputs?: Common_ApplyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_ApplyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_apply(inputs)
	return __es.common_apply(inputs)
});
/**
* | output |
* | --- |
* | "Cancel" |
*
* @param {Common_CancelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const common_cancel = /** @type {((inputs?: Common_CancelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_CancelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_cancel(inputs)
	return __es.common_cancel(inputs)
});
/**
* | output |
* | --- |
* | "Close search" |
*
* @param {Common_Close_SearchInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const common_close_search = /** @type {((inputs?: Common_Close_SearchInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_Close_SearchInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_close_search(inputs)
	return __es.common_close_search(inputs)
});
/**
* | output |
* | --- |
* | "Continue" |
*
* @param {Common_ContinueInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const common_continue = /** @type {((inputs?: Common_ContinueInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_ContinueInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_continue(inputs)
	return __es.common_continue(inputs)
});
/**
* | output |
* | --- |
* | "Value must be a valid number" |
*
* @param {Common_Must_Be_NumberInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const common_must_be_number = /** @type {((inputs?: Common_Must_Be_NumberInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_Must_Be_NumberInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_must_be_number(inputs)
	return __es.common_must_be_number(inputs)
});
/**
* | output |
* | --- |
* | "Preview" |
*
* @param {Common_PreviewInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const common_preview = /** @type {((inputs?: Common_PreviewInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_PreviewInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_preview(inputs)
	return __es.common_preview(inputs)
});
/**
* | output |
* | --- |
* | "Required" |
*
* @param {Common_RequiredInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const common_required = /** @type {((inputs?: Common_RequiredInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_RequiredInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_required(inputs)
	return __es.common_required(inputs)
});
/**
* | output |
* | --- |
* | "Search" |
*
* @param {Common_SearchInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const common_search = /** @type {((inputs?: Common_SearchInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_SearchInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_search(inputs)
	return __es.common_search(inputs)
});
/**
* | output |
* | --- |
* | "Search list" |
*
* @param {Common_Search_ListInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const common_search_list = /** @type {((inputs?: Common_Search_ListInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_Search_ListInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_search_list(inputs)
	return __es.common_search_list(inputs)
});
/**
* | output |
* | --- |
* | "Undo" |
*
* @param {Common_UndoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const common_undo = /** @type {((inputs?: Common_UndoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Common_UndoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.common_undo(inputs)
	return __es.common_undo(inputs)
});
/**
* | output |
* | --- |
* | "Add all dimensions in {name} to rows" |
*
* @param {Dashboard_Add_All_Dims_RowsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_add_all_dims_rows = /** @type {((inputs: Dashboard_Add_All_Dims_RowsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Add_All_Dims_RowsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_add_all_dims_rows(inputs)
	return __es.dashboard_add_all_dims_rows(inputs)
});
/**
* | output |
* | --- |
* | "Add all to columns" |
*
* @param {Dashboard_Add_All_To_ColumnsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_add_all_to_columns = /** @type {((inputs?: Dashboard_Add_All_To_ColumnsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Add_All_To_ColumnsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_add_all_to_columns(inputs)
	return __es.dashboard_add_all_to_columns(inputs)
});
/**
* | output |
* | --- |
* | "Add all in {name} to columns" |
*
* @param {Dashboard_Add_All_To_Columns_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_add_all_to_columns_tag = /** @type {((inputs: Dashboard_Add_All_To_Columns_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Add_All_To_Columns_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_add_all_to_columns_tag(inputs)
	return __es.dashboard_add_all_to_columns_tag(inputs)
});
/**
* | output |
* | --- |
* | "Add all to rows" |
*
* @param {Dashboard_Add_All_To_RowsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_add_all_to_rows = /** @type {((inputs?: Dashboard_Add_All_To_RowsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Add_All_To_RowsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_add_all_to_rows(inputs)
	return __es.dashboard_add_all_to_rows(inputs)
});
/**
* | output |
* | --- |
* | "Add Column" |
*
* @param {Dashboard_Add_ColumnInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_add_column = /** @type {((inputs?: Dashboard_Add_ColumnInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Add_ColumnInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_add_column(inputs)
	return __es.dashboard_add_column(inputs)
});
/**
* | output |
* | --- |
* | "Add filter" |
*
* @param {Dashboard_Add_FilterInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_add_filter = /** @type {((inputs?: Dashboard_Add_FilterInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Add_FilterInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_add_filter(inputs)
	return __es.dashboard_add_filter(inputs)
});
/**
* | output |
* | --- |
* | "Add filter button" |
*
* @param {Dashboard_Add_Filter_ButtonInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_add_filter_button = /** @type {((inputs?: Dashboard_Add_Filter_ButtonInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Add_Filter_ButtonInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_add_filter_button(inputs)
	return __es.dashboard_add_filter_button(inputs)
});
/**
* | output |
* | --- |
* | "Add filter button" |
*
* @param {Dashboard_Add_Filter_Button_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_add_filter_button_aria = /** @type {((inputs?: Dashboard_Add_Filter_Button_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Add_Filter_Button_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_add_filter_button_aria(inputs)
	return __es.dashboard_add_filter_button_aria(inputs)
});
/**
* | output |
* | --- |
* | "Add mock user" |
*
* @param {Dashboard_Add_Mock_UserInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_add_mock_user = /** @type {((inputs?: Dashboard_Add_Mock_UserInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Add_Mock_UserInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_add_mock_user(inputs)
	return __es.dashboard_add_mock_user(inputs)
});
/**
* | output |
* | --- |
* | "Add Row" |
*
* @param {Dashboard_Add_RowInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_add_row = /** @type {((inputs?: Dashboard_Add_RowInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Add_RowInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_add_row(inputs)
	return __es.dashboard_add_row(inputs)
});
/**
* | output |
* | --- |
* | "Add to columns" |
*
* @param {Dashboard_Add_To_ColumnsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_add_to_columns = /** @type {((inputs?: Dashboard_Add_To_ColumnsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Add_To_ColumnsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_add_to_columns(inputs)
	return __es.dashboard_add_to_columns(inputs)
});
/**
* | output |
* | --- |
* | "Add to rows" |
*
* @param {Dashboard_Add_To_RowsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_add_to_rows = /** @type {((inputs?: Dashboard_Add_To_RowsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Add_To_RowsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_add_to_rows(inputs)
	return __es.dashboard_add_to_rows(inputs)
});
/**
* | output |
* | --- |
* | "Added {count} items to filter" |
*
* @param {Dashboard_Added_Items_FilterInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_added_items_filter = /** @type {((inputs: Dashboard_Added_Items_FilterInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Added_Items_FilterInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_added_items_filter(inputs)
	return __es.dashboard_added_items_filter(inputs)
});
/**
* | output |
* | --- |
* | "All" |
*
* @param {Dashboard_AllInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_all = /** @type {((inputs?: Dashboard_AllInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_AllInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_all(inputs)
	return __es.dashboard_all(inputs)
});
/**
* | output |
* | --- |
* | "All measures" |
*
* @param {Dashboard_All_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_all_measures = /** @type {((inputs?: Dashboard_All_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_All_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_all_measures(inputs)
	return __es.dashboard_all_measures(inputs)
});
/**
* | output |
* | --- |
* | "Always show as" |
*
* @param {Dashboard_Always_Show_AsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_always_show_as = /** @type {((inputs?: Dashboard_Always_Show_AsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Always_Show_AsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_always_show_as(inputs)
	return __es.dashboard_always_show_as(inputs)
});
/**
* | output |
* | --- |
* | "Anchor to period end" |
*
* @param {Dashboard_Anchor_Period_EndInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_anchor_period_end = /** @type {((inputs?: Dashboard_Anchor_Period_EndInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Anchor_Period_EndInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_anchor_period_end(inputs)
	return __es.dashboard_anchor_period_end(inputs)
});
/**
* | output |
* | --- |
* | "as of" |
*
* @param {Dashboard_As_OfInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_as_of = /** @type {((inputs?: Dashboard_As_OfInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_As_OfInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_as_of(inputs)
	return __es.dashboard_as_of(inputs)
});
/**
* | output |
* | --- |
* | "as of" |
*
* @param {Dashboard_As_Of_RefInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_as_of_ref = /** @type {((inputs?: Dashboard_As_Of_RefInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_As_Of_RefInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_as_of_ref(inputs)
	return __es.dashboard_as_of_ref(inputs)
});
/**
* | output |
* | --- |
* | "Auto-arrange" |
*
* @param {Dashboard_Auto_ArrangeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_auto_arrange = /** @type {((inputs?: Dashboard_Auto_ArrangeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Auto_ArrangeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_auto_arrange(inputs)
	return __es.dashboard_auto_arrange(inputs)
});
/**
* | output |
* | --- |
* | "Auto-arrange {name}: dimensions to rows, measures to columns" |
*
* @param {Dashboard_Auto_Arrange_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_auto_arrange_tag = /** @type {((inputs: Dashboard_Auto_Arrange_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Auto_Arrange_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_auto_arrange_tag(inputs)
	return __es.dashboard_auto_arrange_tag(inputs)
});
/**
* | output |
* | --- |
* | "Cancel" |
*
* @param {Dashboard_CancelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_cancel = /** @type {((inputs?: Dashboard_CancelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_CancelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_cancel(inputs)
	return __es.dashboard_cancel(inputs)
});
/**
* | output |
* | --- |
* | "Change over comparison period" |
*
* @param {Dashboard_Change_Over_ComparisonInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_change_over_comparison = /** @type {((inputs?: Dashboard_Change_Over_ComparisonInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Change_Over_ComparisonInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_change_over_comparison(inputs)
	return __es.dashboard_change_over_comparison(inputs)
});
/**
* | output |
* | --- |
* | "Clear" |
*
* @param {Dashboard_ClearInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_clear = /** @type {((inputs?: Dashboard_ClearInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_ClearInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_clear(inputs)
	return __es.dashboard_clear(inputs)
});
/**
* | output |
* | --- |
* | "Clear filters" |
*
* @param {Dashboard_Clear_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_clear_filters = /** @type {((inputs?: Dashboard_Clear_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Clear_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_clear_filters(inputs)
	return __es.dashboard_clear_filters(inputs)
});
/**
* | output |
* | --- |
* | "Clear recents" |
*
* @param {Dashboard_Clear_RecentsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_clear_recents = /** @type {((inputs?: Dashboard_Clear_RecentsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Clear_RecentsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_clear_recents(inputs)
	return __es.dashboard_clear_recents(inputs)
});
/**
* | output |
* | --- |
* | "Clear view" |
*
* @param {Dashboard_Clear_ViewInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_clear_view = /** @type {((inputs?: Dashboard_Clear_ViewInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Clear_ViewInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_clear_view(inputs)
	return __es.dashboard_clear_view(inputs)
});
/**
* | output |
* | --- |
* | "Click to edit the values" |
*
* @param {Dashboard_Click_To_Edit_ValuesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_click_to_edit_values = /** @type {((inputs?: Dashboard_Click_To_Edit_ValuesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Click_To_Edit_ValuesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_click_to_edit_values(inputs)
	return __es.dashboard_click_to_edit_values(inputs)
});
/**
* | output |
* | --- |
* | "⌘ + Click to replace" |
*
* @param {Dashboard_Cmd_Click_ReplaceInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_cmd_click_replace = /** @type {((inputs?: Dashboard_Cmd_Click_ReplaceInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Cmd_Click_ReplaceInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_cmd_click_replace(inputs)
	return __es.dashboard_cmd_click_replace(inputs)
});
/**
* | output |
* | --- |
* | "Collapse All" |
*
* @param {Dashboard_Collapse_AllInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_collapse_all = /** @type {((inputs?: Dashboard_Collapse_AllInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Collapse_AllInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_collapse_all(inputs)
	return __es.dashboard_collapse_all(inputs)
});
/**
* | output |
* | --- |
* | "Columns" |
*
* @param {Dashboard_ColumnsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_columns = /** @type {((inputs?: Dashboard_ColumnsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_ColumnsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_columns(inputs)
	return __es.dashboard_columns(inputs)
});
/**
* | output |
* | --- |
* | "Comparison column" |
*
* @param {Dashboard_Comparison_Column_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_comparison_column_aria = /** @type {((inputs?: Dashboard_Comparison_Column_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Comparison_Column_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_comparison_column_aria(inputs)
	return __es.dashboard_comparison_column_aria(inputs)
});
/**
* | output |
* | --- |
* | "Select a comparison for the dashboard" |
*
* @param {Dashboard_Comparison_Select_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_comparison_select_tooltip = /** @type {((inputs?: Dashboard_Comparison_Select_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Comparison_Select_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_comparison_select_tooltip(inputs)
	return __es.dashboard_comparison_select_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "complete data" |
*
* @param {Dashboard_Complete_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_complete_data = /** @type {((inputs?: Dashboard_Complete_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Complete_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_complete_data(inputs)
	return __es.dashboard_complete_data(inputs)
});
/**
* | output |
* | --- |
* | "Timestamp prior to which data frames are considered complete, also known as the watermark" |
*
* @param {Dashboard_Complete_Data_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_complete_data_description = /** @type {((inputs?: Dashboard_Complete_Data_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Complete_Data_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_complete_data_description(inputs)
	return __es.dashboard_complete_data_description(inputs)
});
/**
* | output |
* | --- |
* | "Connect sparse data" |
*
* @param {Dashboard_Connect_Sparse_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_connect_sparse_data = /** @type {((inputs?: Dashboard_Connect_Sparse_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Connect_Sparse_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_connect_sparse_data(inputs)
	return __es.dashboard_connect_sparse_data(inputs)
});
/**
* | output |
* | --- |
* | "Contains" |
*
* @param {Dashboard_ContainsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_contains = /** @type {((inputs?: Dashboard_ContainsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_ContainsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_contains(inputs)
	return __es.dashboard_contains(inputs)
});
/**
* | output |
* | --- |
* | "Copied error to clipboard" |
*
* @param {Dashboard_Copied_Error_ClipboardInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_copied_error_clipboard = /** @type {((inputs?: Dashboard_Copied_Error_ClipboardInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Copied_Error_ClipboardInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_copied_error_clipboard(inputs)
	return __es.dashboard_copied_error_clipboard(inputs)
});
/**
* | output |
* | --- |
* | "Copy full error message" |
*
* @param {Dashboard_Copy_ErrorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_copy_error = /** @type {((inputs?: Dashboard_Copy_ErrorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Copy_ErrorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_copy_error(inputs)
	return __es.dashboard_copy_error(inputs)
});
/**
* | output |
* | --- |
* | "Copy error to clipboard" |
*
* @param {Dashboard_Copy_Error_ClipboardInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_copy_error_clipboard = /** @type {((inputs?: Dashboard_Copy_Error_ClipboardInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Copy_Error_ClipboardInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_copy_error_clipboard(inputs)
	return __es.dashboard_copy_error_clipboard(inputs)
});
/**
* | output |
* | --- |
* | "Copy error message to clipboard" |
*
* @param {Dashboard_Copy_Error_To_Clipboard_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_copy_error_to_clipboard_aria = /** @type {((inputs?: Dashboard_Copy_Error_To_Clipboard_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Copy_Error_To_Clipboard_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_copy_error_to_clipboard_aria(inputs)
	return __es.dashboard_copy_error_to_clipboard_aria(inputs)
});
/**
* | output |
* | --- |
* | "current time" |
*
* @param {Dashboard_Current_TimeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_current_time = /** @type {((inputs?: Dashboard_Current_TimeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Current_TimeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_current_time(inputs)
	return __es.dashboard_current_time(inputs)
});
/**
* | output |
* | --- |
* | "Server clock in selected timezone" |
*
* @param {Dashboard_Current_Time_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_current_time_description = /** @type {((inputs?: Dashboard_Current_Time_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Current_Time_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_current_time_description(inputs)
	return __es.dashboard_current_time_description(inputs)
});
/**
* | output |
* | --- |
* | "Currently showing flat view" |
*
* @param {Dashboard_Currently_FlatInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_currently_flat = /** @type {((inputs?: Dashboard_Currently_FlatInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Currently_FlatInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_currently_flat(inputs)
	return __es.dashboard_currently_flat(inputs)
});
/**
* | output |
* | --- |
* | "Currently showing pivot view" |
*
* @param {Dashboard_Currently_PivotInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_currently_pivot = /** @type {((inputs?: Dashboard_Currently_PivotInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Currently_PivotInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_currently_pivot(inputs)
	return __es.dashboard_currently_pivot(inputs)
});
/**
* | output |
* | --- |
* | "Custom" |
*
* @param {Dashboard_CustomInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_custom = /** @type {((inputs?: Dashboard_CustomInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_CustomInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_custom(inputs)
	return __es.dashboard_custom(inputs)
});
/**
* | output |
* | --- |
* | "Deselect all" |
*
* @param {Dashboard_Deselect_AllInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_deselect_all = /** @type {((inputs?: Dashboard_Deselect_AllInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Deselect_AllInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_deselect_all(inputs)
	return __es.dashboard_deselect_all(inputs)
});
/**
* | output |
* | --- |
* | "Deselect all selections" |
*
* @param {Dashboard_Deselect_All_SelectionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_deselect_all_selections = /** @type {((inputs?: Dashboard_Deselect_All_SelectionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Deselect_All_SelectionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_deselect_all_selections(inputs)
	return __es.dashboard_deselect_all_selections(inputs)
});
/**
* | output |
* | --- |
* | "Dimension Display" |
*
* @param {Dashboard_Dimension_Display_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_dimension_display_aria = /** @type {((inputs?: Dashboard_Dimension_Display_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Dimension_Display_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_dimension_display_aria(inputs)
	return __es.dashboard_dimension_display_aria(inputs)
});
/**
* | output |
* | --- |
* | "Dimension search results" |
*
* @param {Dashboard_Dimension_Search_Results_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_dimension_search_results_aria = /** @type {((inputs?: Dashboard_Dimension_Search_Results_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Dimension_Search_Results_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_dimension_search_results_aria(inputs)
	return __es.dashboard_dimension_search_results_aria(inputs)
});
/**
* | output |
* | --- |
* | "Dimension table" |
*
* @param {Dashboard_Dimension_Table_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_dimension_table_aria = /** @type {((inputs?: Dashboard_Dimension_Table_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Dimension_Table_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_dimension_table_aria(inputs)
	return __es.dashboard_dimension_table_aria(inputs)
});
/**
* | output |
* | --- |
* | "Dimensions" |
*
* @param {Dashboard_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_dimensions = /** @type {((inputs?: Dashboard_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_dimensions(inputs)
	return __es.dashboard_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "{count} dimensions" |
*
* @param {Dashboard_Dimensions_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_dimensions_count = /** @type {((inputs: Dashboard_Dimensions_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Dimensions_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_dimensions_count(inputs)
	return __es.dashboard_dimensions_count(inputs)
});
/**
* | output |
* | --- |
* | "DIMENSIONS" |
*
* @param {Dashboard_Dimensions_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_dimensions_label = /** @type {((inputs?: Dashboard_Dimensions_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Dimensions_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_dimensions_label(inputs)
	return __es.dashboard_dimensions_label(inputs)
});
/**
* | output |
* | --- |
* | "Download PNG" |
*
* @param {Dashboard_Download_PngInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_download_png = /** @type {((inputs?: Dashboard_Download_PngInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Download_PngInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_download_png(inputs)
	return __es.dashboard_download_png(inputs)
});
/**
* | output |
* | --- |
* | "Drag dimensions here" |
*
* @param {Dashboard_Drag_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_drag_dimensions = /** @type {((inputs?: Dashboard_Drag_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Drag_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_drag_dimensions(inputs)
	return __es.dashboard_drag_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "Drag dimensions or measures here" |
*
* @param {Dashboard_Drag_Dimensions_Or_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_drag_dimensions_or_measures = /** @type {((inputs?: Dashboard_Drag_Dimensions_Or_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Drag_Dimensions_Or_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_drag_dimensions_or_measures(inputs)
	return __es.dashboard_drag_dimensions_or_measures(inputs)
});
/**
* | output |
* | --- |
* | "Drag list {zone}" |
*
* @param {Dashboard_Drag_List_ZoneInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_drag_list_zone = /** @type {((inputs: Dashboard_Drag_List_ZoneInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Drag_List_ZoneInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_drag_list_zone(inputs)
	return __es.dashboard_drag_list_zone(inputs)
});
/**
* | output |
* | --- |
* | "Dynamic Y-axis scale" |
*
* @param {Dashboard_Dynamic_Y_AxisInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_dynamic_y_axis = /** @type {((inputs?: Dashboard_Dynamic_Y_AxisInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Dynamic_Y_AxisInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_dynamic_y_axis(inputs)
	return __es.dashboard_dynamic_y_axis(inputs)
});
/**
* | output |
* | --- |
* | "end" |
*
* @param {Dashboard_EndInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_end = /** @type {((inputs?: Dashboard_EndInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_EndInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_end(inputs)
	return __es.dashboard_end(inputs)
});
/**
* | output |
* | --- |
* | "An error occurred" |
*
* @param {Dashboard_Error_OccurredInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_error_occurred = /** @type {((inputs?: Dashboard_Error_OccurredInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Error_OccurredInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_error_occurred(inputs)
	return __es.dashboard_error_occurred(inputs)
});
/**
* | output |
* | --- |
* | "An error occurred. Hover for details." |
*
* @param {Dashboard_Error_Occurred_HoverInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_error_occurred_hover = /** @type {((inputs?: Dashboard_Error_Occurred_HoverInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Error_Occurred_HoverInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_error_occurred_hover(inputs)
	return __es.dashboard_error_occurred_hover(inputs)
});
/**
* | output |
* | --- |
* | "Error" |
*
* @param {Dashboard_Error_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_error_tag = /** @type {((inputs?: Dashboard_Error_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Error_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_error_tag(inputs)
	return __es.dashboard_error_tag(inputs)
});
/**
* | output |
* | --- |
* | "Contact your project's admin for help." |
*
* @param {Dashboard_Errored_Contact_AdminInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_errored_contact_admin = /** @type {((inputs?: Dashboard_Errored_Contact_AdminInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Errored_Contact_AdminInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_errored_contact_admin(inputs)
	return __es.dashboard_errored_contact_admin(inputs)
});
/**
* | output |
* | --- |
* | "Need help? Reach out to us on {link}" |
*
* @param {Dashboard_Errored_Need_HelpInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_errored_need_help = /** @type {((inputs: Dashboard_Errored_Need_HelpInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Errored_Need_HelpInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_errored_need_help(inputs)
	return __es.dashboard_errored_need_help(inputs)
});
/**
* | output |
* | --- |
* | "Sorry, your dashboard isn't working right now!" |
*
* @param {Dashboard_Errored_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_errored_title = /** @type {((inputs?: Dashboard_Errored_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Errored_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_errored_title(inputs)
	return __es.dashboard_errored_title(inputs)
});
/**
* | output |
* | --- |
* | "View project" |
*
* @param {Dashboard_Errored_View_ProjectInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_errored_view_project = /** @type {((inputs?: Dashboard_Errored_View_ProjectInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Errored_View_ProjectInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_errored_view_project(inputs)
	return __es.dashboard_errored_view_project(inputs)
});
/**
* | output |
* | --- |
* | "View project status for errors that may help you find a fix." |
*
* @param {Dashboard_Errored_View_StatusInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_errored_view_status = /** @type {((inputs?: Dashboard_Errored_View_StatusInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Errored_View_StatusInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_errored_view_status(inputs)
	return __es.dashboard_errored_view_status(inputs)
});
/**
* | output |
* | --- |
* | "View project status" |
*
* @param {Dashboard_Errored_View_Status_ButtonInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_errored_view_status_button = /** @type {((inputs?: Dashboard_Errored_View_Status_ButtonInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Errored_View_Status_ButtonInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_errored_view_status_button(inputs)
	return __es.dashboard_errored_view_status_button(inputs)
});
/**
* | output |
* | --- |
* | "Exclude" |
*
* @param {Dashboard_ExcludeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_exclude = /** @type {((inputs?: Dashboard_ExcludeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_ExcludeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_exclude(inputs)
	return __es.dashboard_exclude(inputs)
});
/**
* | output |
* | --- |
* | "Explore" |
*
* @param {Dashboard_ExploreInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_explore = /** @type {((inputs?: Dashboard_ExploreInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_ExploreInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_explore(inputs)
	return __es.dashboard_explore(inputs)
});
/**
* | output |
* | --- |
* | "Export chart" |
*
* @param {Dashboard_Export_ChartInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_export_chart = /** @type {((inputs?: Dashboard_Export_ChartInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Export_ChartInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_export_chart(inputs)
	return __es.dashboard_export_chart(inputs)
});
/**
* | output |
* | --- |
* | "Export dimension table data" |
*
* @param {Dashboard_Export_Dimension_Table_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_export_dimension_table_data = /** @type {((inputs?: Dashboard_Export_Dimension_Table_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Export_Dimension_Table_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_export_dimension_table_data(inputs)
	return __es.dashboard_export_dimension_table_data(inputs)
});
/**
* | output |
* | --- |
* | "Export model data" |
*
* @param {Dashboard_Export_Model_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_export_model_data = /** @type {((inputs?: Dashboard_Export_Model_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Export_Model_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_export_model_data(inputs)
	return __es.dashboard_export_model_data(inputs)
});
/**
* | output |
* | --- |
* | "Export pivot data" |
*
* @param {Dashboard_Export_Pivot_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_export_pivot_data = /** @type {((inputs?: Dashboard_Export_Pivot_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Export_Pivot_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_export_pivot_data(inputs)
	return __es.dashboard_export_pivot_data(inputs)
});
/**
* | output |
* | --- |
* | "Export table data" |
*
* @param {Dashboard_Export_Table_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_export_table_data = /** @type {((inputs?: Dashboard_Export_Table_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Export_Table_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_export_table_data(inputs)
	return __es.dashboard_export_table_data(inputs)
});
/**
* | output |
* | --- |
* | "Filter by this value" |
*
* @param {Dashboard_Filter_By_ValueInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_filter_by_value = /** @type {((inputs?: Dashboard_Filter_By_ValueInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Filter_By_ValueInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_filter_by_value(inputs)
	return __es.dashboard_filter_by_value(inputs)
});
/**
* | output |
* | --- |
* | "Filter dimension value" |
*
* @param {Dashboard_Filter_Dimension_ValueInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_filter_dimension_value = /** @type {((inputs?: Dashboard_Filter_Dimension_ValueInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Filter_Dimension_ValueInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_filter_dimension_value(inputs)
	return __es.dashboard_filter_dimension_value(inputs)
});
/**
* | output |
* | --- |
* | "This filter is required. Set a value to load the dashboard." |
*
* @param {Dashboard_Filter_Required_Set_ValueInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_filter_required_set_value = /** @type {((inputs?: Dashboard_Filter_Required_Set_ValueInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Filter_Required_Set_ValueInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_filter_required_set_value(inputs)
	return __es.dashboard_filter_required_set_value(inputs)
});
/**
* | output |
* | --- |
* | "Flat" |
*
* @param {Dashboard_FlatInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_flat = /** @type {((inputs?: Dashboard_FlatInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_FlatInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_flat(inputs)
	return __es.dashboard_flat(inputs)
});
/**
* | output |
* | --- |
* | "Generated {time}" |
*
* @param {Dashboard_GeneratedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_generated = /** @type {((inputs: Dashboard_GeneratedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_GeneratedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_generated(inputs)
	return __es.dashboard_generated(inputs)
});
/**
* | output |
* | --- |
* | "Generating…" |
*
* @param {Dashboard_GeneratingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_generating = /** @type {((inputs?: Dashboard_GeneratingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_GeneratingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_generating(inputs)
	return __es.dashboard_generating(inputs)
});
/**
* | output |
* | --- |
* | "Grain" |
*
* @param {Dashboard_GrainInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_grain = /** @type {((inputs?: Dashboard_GrainInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_GrainInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_grain(inputs)
	return __es.dashboard_grain(inputs)
});
/**
* | output |
* | --- |
* | "Hide panels" |
*
* @param {Dashboard_Hide_PanelsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_hide_panels = /** @type {((inputs?: Dashboard_Hide_PanelsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Hide_PanelsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_hide_panels(inputs)
	return __es.dashboard_hide_panels(inputs)
});
/**
* | output |
* | --- |
* | "In list" |
*
* @param {Dashboard_In_ListInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_in_list = /** @type {((inputs?: Dashboard_In_ListInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_In_ListInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_in_list(inputs)
	return __es.dashboard_in_list(inputs)
});
/**
* | output |
* | --- |
* | "Include exclude toggle" |
*
* @param {Dashboard_Include_Exclude_ToggleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_include_exclude_toggle = /** @type {((inputs?: Dashboard_Include_Exclude_ToggleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Include_Exclude_ToggleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_include_exclude_toggle(inputs)
	return __es.dashboard_include_exclude_toggle(inputs)
});
/**
* | output |
* | --- |
* | "Invalid time range" |
*
* @param {Dashboard_Invalid_Time_RangeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_invalid_time_range = /** @type {((inputs?: Dashboard_Invalid_Time_RangeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Invalid_Time_RangeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_invalid_time_range(inputs)
	return __es.dashboard_invalid_time_range(inputs)
});
/**
* | output |
* | --- |
* | "Last refreshed {time}" |
*
* @param {Dashboard_Last_Refreshed_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_last_refreshed_ago = /** @type {((inputs: Dashboard_Last_Refreshed_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Last_Refreshed_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_last_refreshed_ago(inputs)
	return __es.dashboard_last_refreshed_ago(inputs)
});
/**
* | output |
* | --- |
* | "latest data" |
*
* @param {Dashboard_Latest_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_latest_data = /** @type {((inputs?: Dashboard_Latest_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Latest_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_latest_data(inputs)
	return __es.dashboard_latest_data(inputs)
});
/**
* | output |
* | --- |
* | "Timestamp of latest data point" |
*
* @param {Dashboard_Latest_Data_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_latest_data_description = /** @type {((inputs?: Dashboard_Latest_Data_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Latest_Data_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_latest_data_description(inputs)
	return __es.dashboard_latest_data_description(inputs)
});
/**
* | output |
* | --- |
* | "Leaderboards" |
*
* @param {Dashboard_Leaderboards_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_leaderboards_aria = /** @type {((inputs?: Dashboard_Leaderboards_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Leaderboards_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_leaderboards_aria(inputs)
	return __es.dashboard_leaderboards_aria(inputs)
});
/**
* | output |
* | --- |
* | "Create a dashboard" |
*
* @param {Dashboard_List_CreateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_list_create = /** @type {((inputs?: Dashboard_List_CreateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_List_CreateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_list_create(inputs)
	return __es.dashboard_list_create(inputs)
});
/**
* | output |
* | --- |
* | "{link} to get started" |
*
* @param {Dashboard_List_Create_To_StartInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_list_create_to_start = /** @type {((inputs: Dashboard_List_Create_To_StartInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_List_Create_To_StartInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_list_create_to_start(inputs)
	return __es.dashboard_list_create_to_start(inputs)
});
/**
* | output |
* | --- |
* | "You don't have any dashboards yet" |
*
* @param {Dashboard_List_EmptyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_list_empty = /** @type {((inputs?: Dashboard_List_EmptyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_List_EmptyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_list_empty(inputs)
	return __es.dashboard_list_empty(inputs)
});
/**
* | output |
* | --- |
* | "See all dashboards" |
*
* @param {Dashboard_List_See_AllInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_list_see_all = /** @type {((inputs?: Dashboard_List_See_AllInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_List_See_AllInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_list_see_all(inputs)
	return __es.dashboard_list_see_all(inputs)
});
/**
* | output |
* | --- |
* | "Measure Chart for {name}" |
*
* @param {Dashboard_Measure_Chart_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_measure_chart_aria = /** @type {((inputs: Dashboard_Measure_Chart_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Measure_Chart_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_measure_chart_aria(inputs)
	return __es.dashboard_measure_chart_aria(inputs)
});
/**
* | output |
* | --- |
* | "Measures" |
*
* @param {Dashboard_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_measures = /** @type {((inputs?: Dashboard_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_measures(inputs)
	return __es.dashboard_measures(inputs)
});
/**
* | output |
* | --- |
* | "{count} measures" |
*
* @param {Dashboard_Measures_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_measures_count = /** @type {((inputs: Dashboard_Measures_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Measures_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_measures_count(inputs)
	return __es.dashboard_measures_count(inputs)
});
/**
* | output |
* | --- |
* | "MEASURES" |
*
* @param {Dashboard_Measures_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_measures_label = /** @type {((inputs?: Dashboard_Measures_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Measures_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_measures_label(inputs)
	return __es.dashboard_measures_label(inputs)
});
/**
* | output |
* | --- |
* | "All Dimensions" |
*
* @param {Dashboard_Menu_All_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_menu_all_dimensions = /** @type {((inputs?: Dashboard_Menu_All_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Menu_All_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_menu_all_dimensions(inputs)
	return __es.dashboard_menu_all_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "No additional details available." |
*
* @param {Dashboard_No_Additional_DetailsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_additional_details = /** @type {((inputs?: Dashboard_No_Additional_DetailsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Additional_DetailsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_additional_details(inputs)
	return __es.dashboard_no_additional_details(inputs)
});
/**
* | output |
* | --- |
* | "No available fields" |
*
* @param {Dashboard_No_Available_FieldsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_available_fields = /** @type {((inputs?: Dashboard_No_Available_FieldsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Available_FieldsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_available_fields(inputs)
	return __es.dashboard_no_available_fields(inputs)
});
/**
* | output |
* | --- |
* | "No comparison dimension selected" |
*
* @param {Dashboard_No_Comparison_DimensionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_comparison_dimension = /** @type {((inputs?: Dashboard_No_Comparison_DimensionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Comparison_DimensionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_comparison_dimension(inputs)
	return __es.dashboard_no_comparison_dimension(inputs)
});
/**
* | output |
* | --- |
* | "No filters selected" |
*
* @param {Dashboard_No_Filters_SelectedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_filters_selected = /** @type {((inputs?: Dashboard_No_Filters_SelectedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Filters_SelectedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_filters_selected(inputs)
	return __es.dashboard_no_filters_selected(inputs)
});
/**
* | output |
* | --- |
* | "No matching tags" |
*
* @param {Dashboard_No_Matching_TagsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_matching_tags = /** @type {((inputs?: Dashboard_No_Matching_TagsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Matching_TagsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_matching_tags(inputs)
	return __es.dashboard_no_matching_tags(inputs)
});
/**
* | output |
* | --- |
* | "No mock users" |
*
* @param {Dashboard_No_Mock_UsersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_mock_users = /** @type {((inputs?: Dashboard_No_Mock_UsersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Mock_UsersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_mock_users(inputs)
	return __es.dashboard_no_mock_users(inputs)
});
/**
* | output |
* | --- |
* | "No options found" |
*
* @param {Dashboard_No_Options_FoundInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_options_found = /** @type {((inputs?: Dashboard_No_Options_FoundInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Options_FoundInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_options_found(inputs)
	return __es.dashboard_no_options_found(inputs)
});
/**
* | output |
* | --- |
* | "No search results to show" |
*
* @param {Dashboard_No_Search_ResultsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_search_results = /** @type {((inputs?: Dashboard_No_Search_ResultsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Search_ResultsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_search_results(inputs)
	return __es.dashboard_no_search_results(inputs)
});
/**
* | output |
* | --- |
* | "No timezones configured" |
*
* @param {Dashboard_No_Timezones_ConfiguredInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_timezones_configured = /** @type {((inputs?: Dashboard_No_Timezones_ConfiguredInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Timezones_ConfiguredInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_timezones_configured(inputs)
	return __es.dashboard_no_timezones_configured(inputs)
});
/**
* | output |
* | --- |
* | "No valid grains available." |
*
* @param {Dashboard_No_Valid_GrainsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_no_valid_grains = /** @type {((inputs?: Dashboard_No_Valid_GrainsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_No_Valid_GrainsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_no_valid_grains(inputs)
	return __es.dashboard_no_valid_grains(inputs)
});
/**
* | output |
* | --- |
* | "of" |
*
* @param {Dashboard_OfInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_of = /** @type {((inputs?: Dashboard_OfInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_OfInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_of(inputs)
	return __es.dashboard_of(inputs)
});
/**
* | output |
* | --- |
* | "Open dimension details" |
*
* @param {Dashboard_Open_Dimension_Details_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_open_dimension_details_aria = /** @type {((inputs?: Dashboard_Open_Dimension_Details_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Open_Dimension_Details_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_open_dimension_details_aria(inputs)
	return __es.dashboard_open_dimension_details_aria(inputs)
});
/**
* | output |
* | --- |
* | "other" |
*
* @param {Dashboard_OtherInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_other = /** @type {((inputs?: Dashboard_OtherInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_OtherInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_other(inputs)
	return __es.dashboard_other(inputs)
});
/**
* | output |
* | --- |
* | "others" |
*
* @param {Dashboard_OthersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_others = /** @type {((inputs?: Dashboard_OthersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_OthersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_others(inputs)
	return __es.dashboard_others(inputs)
});
/**
* | output |
* | --- |
* | "Output excludes selected values" |
*
* @param {Dashboard_Output_ExcludesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_output_excludes = /** @type {((inputs?: Dashboard_Output_ExcludesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Output_ExcludesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_output_excludes(inputs)
	return __es.dashboard_output_excludes(inputs)
});
/**
* | output |
* | --- |
* | "Output includes selected values" |
*
* @param {Dashboard_Output_IncludesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_output_includes = /** @type {((inputs?: Dashboard_Output_IncludesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Output_IncludesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_output_includes(inputs)
	return __es.dashboard_output_includes(inputs)
});
/**
* | output |
* | --- |
* | "Percent of total" |
*
* @param {Dashboard_Percent_Of_TotalInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_percent_of_total = /** @type {((inputs?: Dashboard_Percent_Of_TotalInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Percent_Of_TotalInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_percent_of_total(inputs)
	return __es.dashboard_percent_of_total(inputs)
});
/**
* | output |
* | --- |
* | "Percentage change over comparison period" |
*
* @param {Dashboard_Percentage_ChangeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_percentage_change = /** @type {((inputs?: Dashboard_Percentage_ChangeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Percentage_ChangeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_percentage_change(inputs)
	return __es.dashboard_percentage_change(inputs)
});
/**
* | output |
* | --- |
* | "Pivot" |
*
* @param {Dashboard_PivotInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_pivot = /** @type {((inputs?: Dashboard_PivotInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_PivotInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_pivot(inputs)
	return __es.dashboard_pivot(inputs)
});
/**
* | output |
* | --- |
* | "Add a measure to complete your table." |
*
* @param {Dashboard_Pivot_Add_MeasureInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_pivot_add_measure = /** @type {((inputs?: Dashboard_Pivot_Add_MeasureInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Pivot_Add_MeasureInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_pivot_add_measure(inputs)
	return __es.dashboard_pivot_add_measure(inputs)
});
/**
* | output |
* | --- |
* | "Hang tight! We're building your table..." |
*
* @param {Dashboard_Pivot_Building_TableInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_pivot_building_table = /** @type {((inputs?: Dashboard_Pivot_Building_TableInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Pivot_Building_TableInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_pivot_building_table(inputs)
	return __es.dashboard_pivot_building_table(inputs)
});
/**
* | output |
* | --- |
* | "Give it some data to keep it company." |
*
* @param {Dashboard_Pivot_Give_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_pivot_give_data = /** @type {((inputs?: Dashboard_Pivot_Give_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Pivot_Give_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_pivot_give_data(inputs)
	return __es.dashboard_pivot_give_data(inputs)
});
/**
* | output |
* | --- |
* | "Keep it up!" |
*
* @param {Dashboard_Pivot_Keep_It_UpInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_pivot_keep_it_up = /** @type {((inputs?: Dashboard_Pivot_Keep_It_UpInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Pivot_Keep_It_UpInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_pivot_keep_it_up(inputs)
	return __es.dashboard_pivot_keep_it_up(inputs)
});
/**
* | output |
* | --- |
* | "Learn more about tables in our" |
*
* @param {Dashboard_Pivot_Learn_MoreInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_pivot_learn_more = /** @type {((inputs?: Dashboard_Pivot_Learn_MoreInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Pivot_Learn_MoreInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_pivot_learn_more(inputs)
	return __es.dashboard_pivot_learn_more(inputs)
});
/**
* | output |
* | --- |
* | "Need help? Reach out to us on" |
*
* @param {Dashboard_Pivot_Need_Help_DiscordInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_pivot_need_help_discord = /** @type {((inputs?: Dashboard_Pivot_Need_Help_DiscordInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Pivot_Need_Help_DiscordInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_pivot_need_help_discord(inputs)
	return __es.dashboard_pivot_need_help_discord(inputs)
});
/**
* | output |
* | --- |
* | "No data to show for the selected filters." |
*
* @param {Dashboard_Pivot_No_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_pivot_no_data = /** @type {((inputs?: Dashboard_Pivot_No_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Pivot_No_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_pivot_no_data(inputs)
	return __es.dashboard_pivot_no_data(inputs)
});
/**
* | output |
* | --- |
* | "Your table looks lonely" |
*
* @param {Dashboard_Pivot_Table_LonelyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_pivot_table_lonely = /** @type {((inputs?: Dashboard_Pivot_Table_LonelyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Pivot_Table_LonelyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_pivot_table_lonely(inputs)
	return __es.dashboard_pivot_table_lonely(inputs)
});
/**
* | output |
* | --- |
* | "Readonly Filter Chips" |
*
* @param {Dashboard_Readonly_Filter_Chips_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_readonly_filter_chips_aria = /** @type {((inputs?: Dashboard_Readonly_Filter_Chips_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Readonly_Filter_Chips_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_readonly_filter_chips_aria(inputs)
	return __es.dashboard_readonly_filter_chips_aria(inputs)
});
/**
* | output |
* | --- |
* | "Recent" |
*
* @param {Dashboard_RecentInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_recent = /** @type {((inputs?: Dashboard_RecentInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_RecentInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_recent(inputs)
	return __es.dashboard_recent(inputs)
});
/**
* | output |
* | --- |
* | "Reference" |
*
* @param {Dashboard_ReferenceInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_reference = /** @type {((inputs?: Dashboard_ReferenceInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_ReferenceInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_reference(inputs)
	return __es.dashboard_reference(inputs)
});
/**
* | output |
* | --- |
* | "Remove {label}" |
*
* @param {Dashboard_Remove_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_remove_label = /** @type {((inputs: Dashboard_Remove_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Remove_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_remove_label(inputs)
	return __es.dashboard_remove_label(inputs)
});
/**
* | output |
* | --- |
* | "Removed {count} items from filter" |
*
* @param {Dashboard_Removed_Items_FilterInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_removed_items_filter = /** @type {((inputs: Dashboard_Removed_Items_FilterInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Removed_Items_FilterInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_removed_items_filter(inputs)
	return __es.dashboard_removed_items_filter(inputs)
});
/**
* | output |
* | --- |
* | "Replace" |
*
* @param {Dashboard_ReplaceInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_replace = /** @type {((inputs?: Dashboard_ReplaceInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_ReplaceInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_replace(inputs)
	return __es.dashboard_replace(inputs)
});
/**
* | output |
* | --- |
* | "Replace rows and columns with auto-arranged {name}" |
*
* @param {Dashboard_Replace_Auto_ArrangeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_replace_auto_arrange = /** @type {((inputs: Dashboard_Replace_Auto_ArrangeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Replace_Auto_ArrangeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_replace_auto_arrange(inputs)
	return __es.dashboard_replace_auto_arrange(inputs)
});
/**
* | output |
* | --- |
* | "Replace columns with items in {name}" |
*
* @param {Dashboard_Replace_Columns_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_replace_columns_tag = /** @type {((inputs: Dashboard_Replace_Columns_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Replace_Columns_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_replace_columns_tag(inputs)
	return __es.dashboard_replace_columns_tag(inputs)
});
/**
* | output |
* | --- |
* | "Replace columns with this tag's items" |
*
* @param {Dashboard_Replace_Columns_Tag_ItemsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_replace_columns_tag_items = /** @type {((inputs?: Dashboard_Replace_Columns_Tag_ItemsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Replace_Columns_Tag_ItemsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_replace_columns_tag_items(inputs)
	return __es.dashboard_replace_columns_tag_items(inputs)
});
/**
* | output |
* | --- |
* | "Starting a new table will lose your previous work. Bookmark tables you want to keep" |
*
* @param {Dashboard_Replace_Pivot_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_replace_pivot_description = /** @type {((inputs?: Dashboard_Replace_Pivot_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Replace_Pivot_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_replace_pivot_description(inputs)
	return __es.dashboard_replace_pivot_description(inputs)
});
/**
* | output |
* | --- |
* | "Replace current pivot table?" |
*
* @param {Dashboard_Replace_Pivot_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_replace_pivot_title = /** @type {((inputs?: Dashboard_Replace_Pivot_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Replace_Pivot_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_replace_pivot_title(inputs)
	return __es.dashboard_replace_pivot_title(inputs)
});
/**
* | output |
* | --- |
* | "Replace rows and columns with this tag" |
*
* @param {Dashboard_Replace_Rows_Cols_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_replace_rows_cols_tag = /** @type {((inputs?: Dashboard_Replace_Rows_Cols_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Replace_Rows_Cols_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_replace_rows_cols_tag(inputs)
	return __es.dashboard_replace_rows_cols_tag(inputs)
});
/**
* | output |
* | --- |
* | "Replace rows with dimensions in {name}" |
*
* @param {Dashboard_Replace_Rows_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_replace_rows_tag = /** @type {((inputs: Dashboard_Replace_Rows_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Replace_Rows_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_replace_rows_tag(inputs)
	return __es.dashboard_replace_rows_tag(inputs)
});
/**
* | output |
* | --- |
* | "Replace rows with this tag's dimensions" |
*
* @param {Dashboard_Replace_Rows_Tag_DimsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_replace_rows_tag_dims = /** @type {((inputs?: Dashboard_Replace_Rows_Tag_DimsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Replace_Rows_Tag_DimsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_replace_rows_tag_dims(inputs)
	return __es.dashboard_replace_rows_tag_dims(inputs)
});
/**
* | output |
* | --- |
* | "required measure" |
*
* @param {Dashboard_Required_MeasureInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_required_measure = /** @type {((inputs?: Dashboard_Required_MeasureInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Required_MeasureInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_required_measure(inputs)
	return __es.dashboard_required_measure(inputs)
});
/**
* | output |
* | --- |
* | "Row limit" |
*
* @param {Dashboard_Row_LimitInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_row_limit = /** @type {((inputs?: Dashboard_Row_LimitInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Row_LimitInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_row_limit(inputs)
	return __es.dashboard_row_limit(inputs)
});
/**
* | output |
* | --- |
* | "Only up to top N child rows are shown under each dimension" |
*
* @param {Dashboard_Row_Limit_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_row_limit_tooltip = /** @type {((inputs?: Dashboard_Row_Limit_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Row_Limit_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_row_limit_tooltip(inputs)
	return __es.dashboard_row_limit_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "Rows" |
*
* @param {Dashboard_RowsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_rows = /** @type {((inputs?: Dashboard_RowsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_RowsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_rows(inputs)
	return __es.dashboard_rows(inputs)
});
/**
* | output |
* | --- |
* | "Search Dimension" |
*
* @param {Dashboard_Search_DimensionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_search_dimension = /** @type {((inputs?: Dashboard_Search_DimensionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Search_DimensionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_search_dimension(inputs)
	return __es.dashboard_search_dimension(inputs)
});
/**
* | output |
* | --- |
* | "Search dimensions" |
*
* @param {Dashboard_Search_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_search_dimensions = /** @type {((inputs?: Dashboard_Search_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Search_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_search_dimensions(inputs)
	return __es.dashboard_search_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "Search Results" |
*
* @param {Dashboard_Search_ResultsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_search_results = /** @type {((inputs?: Dashboard_Search_ResultsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Search_ResultsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_search_results(inputs)
	return __es.dashboard_search_results(inputs)
});
/**
* | output |
* | --- |
* | "See more" |
*
* @param {Dashboard_See_MoreInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_see_more = /** @type {((inputs?: Dashboard_See_MoreInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_See_MoreInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_see_more(inputs)
	return __es.dashboard_see_more(inputs)
});
/**
* | output |
* | --- |
* | "Select aggregation grain" |
*
* @param {Dashboard_Select_Aggregation_Grain_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_select_aggregation_grain_aria = /** @type {((inputs?: Dashboard_Select_Aggregation_Grain_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Select_Aggregation_Grain_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_select_aggregation_grain_aria(inputs)
	return __es.dashboard_select_aggregation_grain_aria(inputs)
});
/**
* | output |
* | --- |
* | "Select all" |
*
* @param {Dashboard_Select_AllInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_select_all = /** @type {((inputs?: Dashboard_Select_AllInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Select_AllInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_select_all(inputs)
	return __es.dashboard_select_all(inputs)
});
/**
* | output |
* | --- |
* | "Select a comparison dimension" |
*
* @param {Dashboard_Select_Comparison_DimensionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_select_comparison_dimension = /** @type {((inputs?: Dashboard_Select_Comparison_DimensionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Select_Comparison_DimensionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_select_comparison_dimension(inputs)
	return __es.dashboard_select_comparison_dimension(inputs)
});
/**
* | output |
* | --- |
* | "To see more values, select a comparison dimension above." |
*
* @param {Dashboard_Select_Comparison_HintInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_select_comparison_hint = /** @type {((inputs?: Dashboard_Select_Comparison_HintInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Select_Comparison_HintInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_select_comparison_hint(inputs)
	return __es.dashboard_select_comparison_hint(inputs)
});
/**
* | output |
* | --- |
* | "Select reference time and grain" |
*
* @param {Dashboard_Select_Ref_Time_GrainInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_select_ref_time_grain = /** @type {((inputs?: Dashboard_Select_Ref_Time_GrainInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Select_Ref_Time_GrainInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_select_ref_time_grain(inputs)
	return __es.dashboard_select_ref_time_grain(inputs)
});
/**
* | output |
* | --- |
* | "Select time axis" |
*
* @param {Dashboard_Select_Time_AxisInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_select_time_axis = /** @type {((inputs?: Dashboard_Select_Time_AxisInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Select_Time_AxisInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_select_time_axis(inputs)
	return __es.dashboard_select_time_axis(inputs)
});
/**
* | output |
* | --- |
* | "Select time comparison option" |
*
* @param {Dashboard_Select_Time_Comparison_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_select_time_comparison_aria = /** @type {((inputs?: Dashboard_Select_Time_Comparison_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Select_Time_Comparison_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_select_time_comparison_aria(inputs)
	return __es.dashboard_select_time_comparison_aria(inputs)
});
/**
* | output |
* | --- |
* | "Select {label} time dimension" |
*
* @param {Dashboard_Select_Time_DimensionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_select_time_dimension = /** @type {((inputs: Dashboard_Select_Time_DimensionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Select_Time_DimensionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_select_time_dimension(inputs)
	return __es.dashboard_select_time_dimension(inputs)
});
/**
* | output |
* | --- |
* | "Select a time grain" |
*
* @param {Dashboard_Select_Time_GrainInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_select_time_grain = /** @type {((inputs?: Dashboard_Select_Time_GrainInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Select_Time_GrainInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_select_time_grain(inputs)
	return __es.dashboard_select_time_grain(inputs)
});
/**
* | output |
* | --- |
* | "Select time range" |
*
* @param {Dashboard_Select_Time_RangeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_select_time_range = /** @type {((inputs?: Dashboard_Select_Time_RangeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Select_Time_RangeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_select_time_range(inputs)
	return __es.dashboard_select_time_range(inputs)
});
/**
* | output |
* | --- |
* | "Select time range" |
*
* @param {Dashboard_Select_Time_Range_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_select_time_range_aria = /** @type {((inputs?: Dashboard_Select_Time_Range_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Select_Time_Range_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_select_time_range_aria(inputs)
	return __es.dashboard_select_time_range_aria(inputs)
});
/**
* | output |
* | --- |
* | "selected" |
*
* @param {Dashboard_SelectedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_selected = /** @type {((inputs?: Dashboard_SelectedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_SelectedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_selected(inputs)
	return __es.dashboard_selected(inputs)
});
/**
* | output |
* | --- |
* | "Show panels" |
*
* @param {Dashboard_Show_PanelsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_show_panels = /** @type {((inputs?: Dashboard_Show_PanelsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Show_PanelsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_show_panels(inputs)
	return __es.dashboard_show_panels(inputs)
});
/**
* | output |
* | --- |
* | "Toggle sort leaderboards by absolute change" |
*
* @param {Dashboard_Sort_By_Absolute_Change_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_sort_by_absolute_change_aria = /** @type {((inputs?: Dashboard_Sort_By_Absolute_Change_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Sort_By_Absolute_Change_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_sort_by_absolute_change_aria(inputs)
	return __es.dashboard_sort_by_absolute_change_aria(inputs)
});
/**
* | output |
* | --- |
* | "Toggle sort leaderboards by percent change" |
*
* @param {Dashboard_Sort_By_Percent_Change_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_sort_by_percent_change_aria = /** @type {((inputs?: Dashboard_Sort_By_Percent_Change_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Sort_By_Percent_Change_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_sort_by_percent_change_aria(inputs)
	return __es.dashboard_sort_by_percent_change_aria(inputs)
});
/**
* | output |
* | --- |
* | "Toggle sort leaderboards by percent of total" |
*
* @param {Dashboard_Sort_By_Percent_Total_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_sort_by_percent_total_aria = /** @type {((inputs?: Dashboard_Sort_By_Percent_Total_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Sort_By_Percent_Total_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_sort_by_percent_total_aria(inputs)
	return __es.dashboard_sort_by_percent_total_aria(inputs)
});
/**
* | output |
* | --- |
* | "Toggle sort leaderboards by value" |
*
* @param {Dashboard_Sort_By_Value_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_sort_by_value_aria = /** @type {((inputs?: Dashboard_Sort_By_Value_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Sort_By_Value_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_sort_by_value_aria(inputs)
	return __es.dashboard_sort_by_value_aria(inputs)
});
/**
* | output |
* | --- |
* | "start" |
*
* @param {Dashboard_StartInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_start = /** @type {((inputs?: Dashboard_StartInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_StartInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_start(inputs)
	return __es.dashboard_start(inputs)
});
/**
* | output |
* | --- |
* | "Start Pivot" |
*
* @param {Dashboard_Start_PivotInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_start_pivot = /** @type {((inputs?: Dashboard_Start_PivotInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Start_PivotInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_start_pivot(inputs)
	return __es.dashboard_start_pivot(inputs)
});
/**
* | output |
* | --- |
* | "Switch to flat view" |
*
* @param {Dashboard_Switch_FlatInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_switch_flat = /** @type {((inputs?: Dashboard_Switch_FlatInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Switch_FlatInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_switch_flat(inputs)
	return __es.dashboard_switch_flat(inputs)
});
/**
* | output |
* | --- |
* | "Switch to pivot view" |
*
* @param {Dashboard_Switch_PivotInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_switch_pivot = /** @type {((inputs?: Dashboard_Switch_PivotInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Switch_PivotInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_switch_pivot(inputs)
	return __es.dashboard_switch_pivot(inputs)
});
/**
* | output |
* | --- |
* | "Table mode" |
*
* @param {Dashboard_Table_ModeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_table_mode = /** @type {((inputs?: Dashboard_Table_ModeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Table_ModeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_table_mode(inputs)
	return __es.dashboard_table_mode(inputs)
});
/**
* | output |
* | --- |
* | "Tags" |
*
* @param {Dashboard_TagsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_tags = /** @type {((inputs?: Dashboard_TagsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_TagsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_tags(inputs)
	return __es.dashboard_tags(inputs)
});
/**
* | output |
* | --- |
* | "If the issue persists, please contact us on" |
*
* @param {Dashboard_Tdd_Contact_DiscordInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_tdd_contact_discord = /** @type {((inputs?: Dashboard_Tdd_Contact_DiscordInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Tdd_Contact_DiscordInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_tdd_contact_discord(inputs)
	return __es.dashboard_tdd_contact_discord(inputs)
});
/**
* | output |
* | --- |
* | "We encountered an error while loading the data. Please try refreshing the page." |
*
* @param {Dashboard_Tdd_ErrorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_tdd_error = /** @type {((inputs?: Dashboard_Tdd_ErrorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Tdd_ErrorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_tdd_error(inputs)
	return __es.dashboard_tdd_error(inputs)
});
/**
* | output |
* | --- |
* | "No Comparison" |
*
* @param {Dashboard_Tdd_No_ComparisonInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_tdd_no_comparison = /** @type {((inputs?: Dashboard_Tdd_No_ComparisonInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Tdd_No_ComparisonInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_tdd_no_comparison(inputs)
	return __es.dashboard_tdd_no_comparison(inputs)
});
/**
* | output |
* | --- |
* | "Time" |
*
* @param {Dashboard_Tdd_TimeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_tdd_time = /** @type {((inputs?: Dashboard_Tdd_TimeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Tdd_TimeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_tdd_time(inputs)
	return __es.dashboard_tdd_time(inputs)
});
/**
* | output |
* | --- |
* | "Time" |
*
* @param {Dashboard_TimeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_time = /** @type {((inputs?: Dashboard_TimeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_TimeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_time(inputs)
	return __es.dashboard_time(inputs)
});
/**
* | output |
* | --- |
* | "Time axis" |
*
* @param {Dashboard_Time_AxisInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_time_axis = /** @type {((inputs?: Dashboard_Time_AxisInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Time_AxisInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_time_axis(inputs)
	return __es.dashboard_time_axis(inputs)
});
/**
* | output |
* | --- |
* | "TIME" |
*
* @param {Dashboard_Time_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_time_label = /** @type {((inputs?: Dashboard_Time_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Time_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_time_label(inputs)
	return __es.dashboard_time_label(inputs)
});
/**
* | output |
* | --- |
* | "Time zone" |
*
* @param {Dashboard_Time_ZoneInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_time_zone = /** @type {((inputs?: Dashboard_Time_ZoneInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Time_ZoneInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_time_zone(inputs)
	return __es.dashboard_time_zone(inputs)
});
/**
* | output |
* | --- |
* | "Timezone selector" |
*
* @param {Dashboard_Timezone_SelectorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_timezone_selector = /** @type {((inputs?: Dashboard_Timezone_SelectorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Timezone_SelectorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_timezone_selector(inputs)
	return __es.dashboard_timezone_selector(inputs)
});
/**
* | output |
* | --- |
* | "to" |
*
* @param {Dashboard_ToInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_to = /** @type {((inputs?: Dashboard_ToInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_ToInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_to(inputs)
	return __es.dashboard_to(inputs)
});
/**
* | output |
* | --- |
* | "Toggle to exclude values" |
*
* @param {Dashboard_Toggle_ExcludeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_toggle_exclude = /** @type {((inputs?: Dashboard_Toggle_ExcludeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Toggle_ExcludeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_toggle_exclude(inputs)
	return __es.dashboard_toggle_exclude(inputs)
});
/**
* | output |
* | --- |
* | "Toggle to include values" |
*
* @param {Dashboard_Toggle_IncludeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_toggle_include = /** @type {((inputs?: Dashboard_Toggle_IncludeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Toggle_IncludeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_toggle_include(inputs)
	return __es.dashboard_toggle_include(inputs)
});
/**
* | output |
* | --- |
* | "Toggle rows viewer" |
*
* @param {Dashboard_Toggle_Rows_Viewer_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_toggle_rows_viewer_aria = /** @type {((inputs?: Dashboard_Toggle_Rows_Viewer_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Toggle_Rows_Viewer_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_toggle_rows_viewer_aria(inputs)
	return __es.dashboard_toggle_rows_viewer_aria(inputs)
});
/**
* | output |
* | --- |
* | "Toggle time comparison" |
*
* @param {Dashboard_Toggle_Time_Comparison_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_toggle_time_comparison_aria = /** @type {((inputs?: Dashboard_Toggle_Time_Comparison_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Toggle_Time_Comparison_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_toggle_time_comparison_aria(inputs)
	return __es.dashboard_toggle_time_comparison_aria(inputs)
});
/**
* | output |
* | --- |
* | "Total column" |
*
* @param {Dashboard_Total_ColumnInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_total_column = /** @type {((inputs?: Dashboard_Total_ColumnInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Total_ColumnInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_total_column(inputs)
	return __es.dashboard_total_column(inputs)
});
/**
* | output |
* | --- |
* | "Total row" |
*
* @param {Dashboard_Total_RowInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_total_row = /** @type {((inputs?: Dashboard_Total_RowInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Total_RowInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_total_row(inputs)
	return __es.dashboard_total_row(inputs)
});
/**
* | output |
* | --- |
* | "View as" |
*
* @param {Dashboard_View_AsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_view_as = /** @type {((inputs?: Dashboard_View_AsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_View_AsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_view_as(inputs)
	return __es.dashboard_view_as(inputs)
});
/**
* | output |
* | --- |
* | "Viewing as" |
*
* @param {Dashboard_Viewing_AsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dashboard_viewing_as = /** @type {((inputs?: Dashboard_Viewing_AsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dashboard_Viewing_AsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dashboard_viewing_as(inputs)
	return __es.dashboard_viewing_as(inputs)
});
/**
* | output |
* | --- |
* | "You haven't saved changes to this alert yet, so closing this window will lose your work." |
*
* @param {Dialog_Close_Without_Saving_Alert_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dialog_close_without_saving_alert_desc = /** @type {((inputs?: Dialog_Close_Without_Saving_Alert_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dialog_Close_Without_Saving_Alert_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dialog_close_without_saving_alert_desc(inputs)
	return __es.dialog_close_without_saving_alert_desc(inputs)
});
/**
* | output |
* | --- |
* | "Keep editing" |
*
* @param {Dialog_Close_Without_Saving_CancelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dialog_close_without_saving_cancel = /** @type {((inputs?: Dialog_Close_Without_Saving_CancelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dialog_Close_Without_Saving_CancelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dialog_close_without_saving_cancel(inputs)
	return __es.dialog_close_without_saving_cancel(inputs)
});
/**
* | output |
* | --- |
* | "Close" |
*
* @param {Dialog_Close_Without_Saving_ConfirmInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dialog_close_without_saving_confirm = /** @type {((inputs?: Dialog_Close_Without_Saving_ConfirmInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dialog_Close_Without_Saving_ConfirmInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dialog_close_without_saving_confirm(inputs)
	return __es.dialog_close_without_saving_confirm(inputs)
});
/**
* | output |
* | --- |
* | "Close without saving?" |
*
* @param {Dialog_Close_Without_Saving_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const dialog_close_without_saving_title = /** @type {((inputs?: Dialog_Close_Without_Saving_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Dialog_Close_Without_Saving_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.dialog_close_without_saving_title(inputs)
	return __es.dialog_close_without_saving_title(inputs)
});
/**
* | output |
* | --- |
* | "You don't have access to this page. Please check that you have the correct permissions." |
*
* @param {Error_Access_Denied_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_access_denied_body = /** @type {((inputs?: Error_Access_Denied_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Access_Denied_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_access_denied_body(inputs)
	return __es.error_access_denied_body(inputs)
});
/**
* | output |
* | --- |
* | "Access denied" |
*
* @param {Error_Access_Denied_HeaderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_access_denied_header = /** @type {((inputs?: Error_Access_Denied_HeaderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Access_Denied_HeaderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_access_denied_header(inputs)
	return __es.error_access_denied_header(inputs)
});
/**
* | output |
* | --- |
* | "Try refreshing the page. If the problem persists, try signing out and back in." |
*
* @param {Error_Auth_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_auth_body = /** @type {((inputs?: Error_Auth_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Auth_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_auth_body(inputs)
	return __es.error_auth_body(inputs)
});
/**
* | output |
* | --- |
* | "Authentication error" |
*
* @param {Error_Auth_HeaderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_auth_header = /** @type {((inputs?: Error_Auth_HeaderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Auth_HeaderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_auth_header(inputs)
	return __es.error_auth_header(inputs)
});
/**
* | output |
* | --- |
* | "Back to home" |
*
* @param {Error_Back_To_HomeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_back_to_home = /** @type {((inputs?: Error_Back_To_HomeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Back_To_HomeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_back_to_home(inputs)
	return __es.error_back_to_home(inputs)
});
/**
* | output |
* | --- |
* | "Please check that you have the correct link or if you have access to it." |
*
* @param {Error_Conversation_Not_Found_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_conversation_not_found_body = /** @type {((inputs?: Error_Conversation_Not_Found_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Conversation_Not_Found_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_conversation_not_found_body(inputs)
	return __es.error_conversation_not_found_body(inputs)
});
/**
* | output |
* | --- |
* | "Conversation not found" |
*
* @param {Error_Conversation_Not_Found_HeaderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_conversation_not_found_header = /** @type {((inputs?: Error_Conversation_Not_Found_HeaderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Conversation_Not_Found_HeaderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_conversation_not_found_header(inputs)
	return __es.error_conversation_not_found_header(inputs)
});
/**
* | output |
* | --- |
* | "There was an error deploying your project. Please contact support." |
*
* @param {Error_Deploying_ProjectInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_deploying_project = /** @type {((inputs?: Error_Deploying_ProjectInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Deploying_ProjectInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_deploying_project(inputs)
	return __es.error_deploying_project(inputs)
});
/**
* | output |
* | --- |
* | "Deployment Error" |
*
* @param {Error_Deployment_ErrorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_deployment_error = /** @type {((inputs?: Error_Deployment_ErrorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Deployment_ErrorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_deployment_error(inputs)
	return __es.error_deployment_error(inputs)
});
/**
* | output |
* | --- |
* | "This is potentially a temporary state if the project has just been reset." |
*
* @param {Error_Deployment_Not_Found_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_deployment_not_found_body = /** @type {((inputs?: Error_Deployment_Not_Found_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Deployment_Not_Found_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_deployment_not_found_body(inputs)
	return __es.error_deployment_not_found_body(inputs)
});
/**
* | output |
* | --- |
* | "Project deployment not found" |
*
* @param {Error_Deployment_Not_Found_HeaderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_deployment_not_found_header = /** @type {((inputs?: Error_Deployment_Not_Found_HeaderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Deployment_Not_Found_HeaderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_deployment_not_found_header(inputs)
	return __es.error_deployment_not_found_header(inputs)
});
/**
* | output |
* | --- |
* | "Error fetching deployment" |
*
* @param {Error_Fetching_DeploymentInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_fetching_deployment = /** @type {((inputs?: Error_Fetching_DeploymentInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Fetching_DeploymentInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_fetching_deployment(inputs)
	return __es.error_fetching_deployment(inputs)
});
/**
* | output |
* | --- |
* | "Try refreshing the page, and reach out to us if the problem persists." |
*
* @param {Error_Generic_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_generic_body = /** @type {((inputs?: Error_Generic_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Generic_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_generic_body(inputs)
	return __es.error_generic_body(inputs)
});
/**
* | output |
* | --- |
* | "Sorry, something went wrong!" |
*
* @param {Error_Generic_HeaderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_generic_header = /** @type {((inputs?: Error_Generic_HeaderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Generic_HeaderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_generic_header(inputs)
	return __es.error_generic_header(inputs)
});
/**
* | output |
* | --- |
* | "Hide details" |
*
* @param {Error_Hide_DetailsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_hide_details = /** @type {((inputs?: Error_Hide_DetailsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Hide_DetailsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_hide_details(inputs)
	return __es.error_hide_details(inputs)
});
/**
* | output |
* | --- |
* | "It looks like this link is no longer active. Please reach out to the sender to request a new link." |
*
* @param {Error_Link_Expired_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_link_expired_body = /** @type {((inputs?: Error_Link_Expired_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Link_Expired_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_link_expired_body(inputs)
	return __es.error_link_expired_body(inputs)
});
/**
* | output |
* | --- |
* | "Oops! This link has expired" |
*
* @param {Error_Link_Expired_HeaderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_link_expired_header = /** @type {((inputs?: Error_Link_Expired_HeaderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Link_Expired_HeaderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_link_expired_header(inputs)
	return __es.error_link_expired_header(inputs)
});
/**
* | output |
* | --- |
* | "It seems we're having trouble reaching our servers. Check your connection or try again later." |
*
* @param {Error_Network_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_network_body = /** @type {((inputs?: Error_Network_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Network_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_network_body(inputs)
	return __es.error_network_body(inputs)
});
/**
* | output |
* | --- |
* | "Network Error" |
*
* @param {Error_Network_HeaderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_network_header = /** @type {((inputs?: Error_Network_HeaderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Network_HeaderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_network_header(inputs)
	return __es.error_network_header(inputs)
});
/**
* | output |
* | --- |
* | "The organization you requested could not be found. Please check that you have provided a valid organization name." |
*
* @param {Error_Org_Not_Found_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_org_not_found_body = /** @type {((inputs?: Error_Org_Not_Found_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Org_Not_Found_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_org_not_found_body(inputs)
	return __es.error_org_not_found_body(inputs)
});
/**
* | output |
* | --- |
* | "Organization not found" |
*
* @param {Error_Org_Not_Found_HeaderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_org_not_found_header = /** @type {((inputs?: Error_Org_Not_Found_HeaderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Org_Not_Found_HeaderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_org_not_found_header(inputs)
	return __es.error_org_not_found_header(inputs)
});
/**
* | output |
* | --- |
* | "The page you're looking for might have been removed, had its name changed, or is temporarily unavailable." |
*
* @param {Error_Page_Not_Found_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_page_not_found_body = /** @type {((inputs?: Error_Page_Not_Found_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Page_Not_Found_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_page_not_found_body(inputs)
	return __es.error_page_not_found_body(inputs)
});
/**
* | output |
* | --- |
* | "Sorry, we can't find this page!" |
*
* @param {Error_Page_Not_Found_HeaderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_page_not_found_header = /** @type {((inputs?: Error_Page_Not_Found_HeaderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Page_Not_Found_HeaderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_page_not_found_header(inputs)
	return __es.error_page_not_found_header(inputs)
});
/**
* | output |
* | --- |
* | "The project you requested could not be found. Please check that you have provided a valid project name." |
*
* @param {Error_Project_Not_Found_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_project_not_found_body = /** @type {((inputs?: Error_Project_Not_Found_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Project_Not_Found_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_project_not_found_body(inputs)
	return __es.error_project_not_found_body(inputs)
});
/**
* | output |
* | --- |
* | "Project not found" |
*
* @param {Error_Project_Not_Found_HeaderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_project_not_found_header = /** @type {((inputs?: Error_Project_Not_Found_HeaderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Project_Not_Found_HeaderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_project_not_found_header(inputs)
	return __es.error_project_not_found_header(inputs)
});
/**
* | output |
* | --- |
* | "This resource may have been deleted, renamed, or is temporarily unavailable." |
*
* @param {Error_Resource_Not_Found_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_resource_not_found_body = /** @type {((inputs?: Error_Resource_Not_Found_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Resource_Not_Found_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_resource_not_found_body(inputs)
	return __es.error_resource_not_found_body(inputs)
});
/**
* | output |
* | --- |
* | "Resource not found" |
*
* @param {Error_Resource_Not_Found_HeaderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_resource_not_found_header = /** @type {((inputs?: Error_Resource_Not_Found_HeaderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Resource_Not_Found_HeaderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_resource_not_found_header(inputs)
	return __es.error_resource_not_found_header(inputs)
});
/**
* | output |
* | --- |
* | "Retry now" |
*
* @param {Error_Retry_NowInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_retry_now = /** @type {((inputs?: Error_Retry_NowInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Retry_NowInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_retry_now(inputs)
	return __es.error_retry_now(inputs)
});
/**
* | output |
* | --- |
* | "Show details" |
*
* @param {Error_Show_DetailsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const error_show_details = /** @type {((inputs?: Error_Show_DetailsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Error_Show_DetailsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.error_show_details(inputs)
	return __es.error_show_details(inputs)
});
/**
* | output |
* | --- |
* | "All Dimensions" |
*
* @param {Explore_All_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_all_dimensions = /** @type {((inputs?: Explore_All_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_All_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_all_dimensions(inputs)
	return __es.explore_all_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "All Measures" |
*
* @param {Explore_All_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_all_measures = /** @type {((inputs?: Explore_All_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_All_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_all_measures(inputs)
	return __es.explore_all_measures(inputs)
});
/**
* | output |
* | --- |
* | "by" |
*
* @param {Explore_By_Grain_PrefixInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_by_grain_prefix = /** @type {((inputs?: Explore_By_Grain_PrefixInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_By_Grain_PrefixInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_by_grain_prefix(inputs)
	return __es.explore_by_grain_prefix(inputs)
});
/**
* | output |
* | --- |
* | "Choose dimensions to display" |
*
* @param {Explore_Choose_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_choose_dimensions = /** @type {((inputs?: Explore_Choose_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Choose_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_choose_dimensions(inputs)
	return __es.explore_choose_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "Choose measures to display" |
*
* @param {Explore_Choose_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_choose_measures = /** @type {((inputs?: Explore_Choose_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Choose_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_choose_measures(inputs)
	return __es.explore_choose_measures(inputs)
});
/**
* | output |
* | --- |
* | "Clear filter" |
*
* @param {Explore_Clear_FilterInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_clear_filter = /** @type {((inputs?: Explore_Clear_FilterInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Clear_FilterInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_clear_filter(inputs)
	return __es.explore_clear_filter(inputs)
});
/**
* | output |
* | --- |
* | "Clear filter {tag}" |
*
* @param {Explore_Clear_Filter_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_clear_filter_tag = /** @type {((inputs: Explore_Clear_Filter_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Clear_Filter_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_clear_filter_tag(inputs)
	return __es.explore_clear_filter_tag(inputs)
});
/**
* | output |
* | --- |
* | "Clear search to reorder dimensions." |
*
* @param {Explore_Clear_Search_To_Reorder_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_clear_search_to_reorder_dimensions = /** @type {((inputs?: Explore_Clear_Search_To_Reorder_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Clear_Search_To_Reorder_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_clear_search_to_reorder_dimensions(inputs)
	return __es.explore_clear_search_to_reorder_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "Clear search to reorder measures." |
*
* @param {Explore_Clear_Search_To_Reorder_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_clear_search_to_reorder_measures = /** @type {((inputs?: Explore_Clear_Search_To_Reorder_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Clear_Search_To_Reorder_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_clear_search_to_reorder_measures(inputs)
	return __es.explore_clear_search_to_reorder_measures(inputs)
});
/**
* | output |
* | --- |
* | "Clear tag filter" |
*
* @param {Explore_Clear_Tag_FilterInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_clear_tag_filter = /** @type {((inputs?: Explore_Clear_Tag_FilterInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Clear_Tag_FilterInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_clear_tag_filter(inputs)
	return __es.explore_clear_tag_filter(inputs)
});
/**
* | output |
* | --- |
* | "Clear the tag filter to reorder dimensions." |
*
* @param {Explore_Clear_Tag_Filter_To_Reorder_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_clear_tag_filter_to_reorder_dimensions = /** @type {((inputs?: Explore_Clear_Tag_Filter_To_Reorder_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Clear_Tag_Filter_To_Reorder_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_clear_tag_filter_to_reorder_dimensions(inputs)
	return __es.explore_clear_tag_filter_to_reorder_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "Clear the tag filter to reorder measures." |
*
* @param {Explore_Clear_Tag_Filter_To_Reorder_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_clear_tag_filter_to_reorder_measures = /** @type {((inputs?: Explore_Clear_Tag_Filter_To_Reorder_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Clear_Tag_Filter_To_Reorder_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_clear_tag_filter_to_reorder_measures(inputs)
	return __es.explore_clear_tag_filter_to_reorder_measures(inputs)
});
/**
* | output |
* | --- |
* | "{count} of {total} Dimensions" |
*
* @param {Explore_Dimensions_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_dimensions_count = /** @type {((inputs: Explore_Dimensions_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Dimensions_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_dimensions_count(inputs)
	return __es.explore_dimensions_count(inputs)
});
/**
* | output |
* | --- |
* | "Filter by {tag}" |
*
* @param {Explore_Filter_By_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_filter_by_tag = /** @type {((inputs: Explore_Filter_By_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Filter_By_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_filter_by_tag(inputs)
	return __es.explore_filter_by_tag(inputs)
});
/**
* | output |
* | --- |
* | "Filter" |
*
* @param {Explore_Filter_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_filter_label = /** @type {((inputs?: Explore_Filter_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Filter_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_filter_label(inputs)
	return __es.explore_filter_label(inputs)
});
/**
* | output |
* | --- |
* | "Go to Explore Dashboard" |
*
* @param {Explore_Go_To_DashboardInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_go_to_dashboard = /** @type {((inputs?: Explore_Go_To_DashboardInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Go_To_DashboardInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_go_to_dashboard(inputs)
	return __es.explore_go_to_dashboard(inputs)
});
/**
* | output |
* | --- |
* | "Go to Explore" |
*
* @param {Explore_Go_To_ExploreInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_go_to_explore = /** @type {((inputs?: Explore_Go_To_ExploreInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Go_To_ExploreInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_go_to_explore(inputs)
	return __es.explore_go_to_explore(inputs)
});
/**
* | output |
* | --- |
* | "Go to {name}" |
*
* @param {Explore_Go_To_NamedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_go_to_named = /** @type {((inputs: Explore_Go_To_NamedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Go_To_NamedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_go_to_named(inputs)
	return __es.explore_go_to_named(inputs)
});
/**
* | output |
* | --- |
* | "Hidden dimensions" |
*
* @param {Explore_Hidden_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_hidden_dimensions = /** @type {((inputs?: Explore_Hidden_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Hidden_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_hidden_dimensions(inputs)
	return __es.explore_hidden_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "Hidden measures" |
*
* @param {Explore_Hidden_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_hidden_measures = /** @type {((inputs?: Explore_Hidden_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Hidden_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_hidden_measures(inputs)
	return __es.explore_hidden_measures(inputs)
});
/**
* | output |
* | --- |
* | "Hide all" |
*
* @param {Explore_Hide_AllInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_hide_all = /** @type {((inputs?: Explore_Hide_AllInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Hide_AllInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_hide_all(inputs)
	return __es.explore_hide_all(inputs)
});
/**
* | output |
* | --- |
* | "Hide all in {tag}" |
*
* @param {Explore_Hide_All_In_Named_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_hide_all_in_named_tag = /** @type {((inputs: Explore_Hide_All_In_Named_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Hide_All_In_Named_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_hide_all_in_named_tag(inputs)
	return __es.explore_hide_all_in_named_tag(inputs)
});
/**
* | output |
* | --- |
* | "Hide all in tag" |
*
* @param {Explore_Hide_All_In_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_hide_all_in_tag = /** @type {((inputs?: Explore_Hide_All_In_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Hide_All_In_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_hide_all_in_tag(inputs)
	return __es.explore_hide_all_in_tag(inputs)
});
/**
* | output |
* | --- |
* | "Hide {name}" |
*
* @param {Explore_Hide_ItemInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_hide_item = /** @type {((inputs: Explore_Hide_ItemInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Hide_ItemInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_hide_item(inputs)
	return __es.explore_hide_item(inputs)
});
/**
* | output |
* | --- |
* | "{count} of {total} Measures" |
*
* @param {Explore_Measures_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_measures_count = /** @type {((inputs: Explore_Measures_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Measures_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_measures_count(inputs)
	return __es.explore_measures_count(inputs)
});
/**
* | output |
* | --- |
* | "Multi-select" |
*
* @param {Explore_Multi_SelectInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_multi_select = /** @type {((inputs?: Explore_Multi_SelectInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Multi_SelectInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_multi_select(inputs)
	return __es.explore_multi_select(inputs)
});
/**
* | output |
* | --- |
* | "Must show at least one dimension" |
*
* @param {Explore_Must_Show_One_DimensionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_must_show_one_dimension = /** @type {((inputs?: Explore_Must_Show_One_DimensionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Must_Show_One_DimensionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_must_show_one_dimension(inputs)
	return __es.explore_must_show_one_dimension(inputs)
});
/**
* | output |
* | --- |
* | "Must show at least one measure" |
*
* @param {Explore_Must_Show_One_MeasureInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_must_show_one_measure = /** @type {((inputs?: Explore_Must_Show_One_MeasureInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Must_Show_One_MeasureInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_must_show_one_measure(inputs)
	return __es.explore_must_show_one_measure(inputs)
});
/**
* | output |
* | --- |
* | "{count} measures" |
*
* @param {Explore_N_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_n_measures = /** @type {((inputs: Explore_N_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_N_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_n_measures(inputs)
	return __es.explore_n_measures(inputs)
});
/**
* | output |
* | --- |
* | "No dimensions from this tag are shown" |
*
* @param {Explore_No_Dimensions_From_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_dimensions_from_tag = /** @type {((inputs?: Explore_No_Dimensions_From_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Dimensions_From_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_dimensions_from_tag(inputs)
	return __es.explore_no_dimensions_from_tag(inputs)
});
/**
* | output |
* | --- |
* | "No dimensions or tags found" |
*
* @param {Explore_No_Dimensions_Or_TagsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_dimensions_or_tags = /** @type {((inputs?: Explore_No_Dimensions_Or_TagsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Dimensions_Or_TagsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_dimensions_or_tags(inputs)
	return __es.explore_no_dimensions_or_tags(inputs)
});
/**
* | output |
* | --- |
* | "No dimensions shown" |
*
* @param {Explore_No_Dimensions_ShownInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_dimensions_shown = /** @type {((inputs?: Explore_No_Dimensions_ShownInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Dimensions_ShownInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_dimensions_shown(inputs)
	return __es.explore_no_dimensions_shown(inputs)
});
/**
* | output |
* | --- |
* | "No hidden dimensions" |
*
* @param {Explore_No_Hidden_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_hidden_dimensions = /** @type {((inputs?: Explore_No_Hidden_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Hidden_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_hidden_dimensions(inputs)
	return __es.explore_no_hidden_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "No hidden measures" |
*
* @param {Explore_No_Hidden_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_hidden_measures = /** @type {((inputs?: Explore_No_Hidden_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Hidden_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_hidden_measures(inputs)
	return __es.explore_no_hidden_measures(inputs)
});
/**
* | output |
* | --- |
* | "No matching dimensions shown" |
*
* @param {Explore_No_Matching_Dimensions_ShownInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_matching_dimensions_shown = /** @type {((inputs?: Explore_No_Matching_Dimensions_ShownInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Matching_Dimensions_ShownInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_matching_dimensions_shown(inputs)
	return __es.explore_no_matching_dimensions_shown(inputs)
});
/**
* | output |
* | --- |
* | "No matching hidden dimensions" |
*
* @param {Explore_No_Matching_Hidden_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_matching_hidden_dimensions = /** @type {((inputs?: Explore_No_Matching_Hidden_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Matching_Hidden_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_matching_hidden_dimensions(inputs)
	return __es.explore_no_matching_hidden_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "No matching hidden measures" |
*
* @param {Explore_No_Matching_Hidden_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_matching_hidden_measures = /** @type {((inputs?: Explore_No_Matching_Hidden_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Matching_Hidden_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_matching_hidden_measures(inputs)
	return __es.explore_no_matching_hidden_measures(inputs)
});
/**
* | output |
* | --- |
* | "No matching leaderboard measures shown" |
*
* @param {Explore_No_Matching_Leaderboard_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_matching_leaderboard_measures = /** @type {((inputs?: Explore_No_Matching_Leaderboard_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Matching_Leaderboard_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_matching_leaderboard_measures(inputs)
	return __es.explore_no_matching_leaderboard_measures(inputs)
});
/**
* | output |
* | --- |
* | "No matching measures shown" |
*
* @param {Explore_No_Matching_Measures_ShownInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_matching_measures_shown = /** @type {((inputs?: Explore_No_Matching_Measures_ShownInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Matching_Measures_ShownInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_matching_measures_shown(inputs)
	return __es.explore_no_matching_measures_shown(inputs)
});
/**
* | output |
* | --- |
* | "No matching tags" |
*
* @param {Explore_No_Matching_TagsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_matching_tags = /** @type {((inputs?: Explore_No_Matching_TagsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Matching_TagsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_matching_tags(inputs)
	return __es.explore_no_matching_tags(inputs)
});
/**
* | output |
* | --- |
* | "No measures from this tag are shown" |
*
* @param {Explore_No_Measures_From_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_measures_from_tag = /** @type {((inputs?: Explore_No_Measures_From_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Measures_From_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_measures_from_tag(inputs)
	return __es.explore_no_measures_from_tag(inputs)
});
/**
* | output |
* | --- |
* | "No measures or tags found" |
*
* @param {Explore_No_Measures_Or_TagsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_measures_or_tags = /** @type {((inputs?: Explore_No_Measures_Or_TagsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Measures_Or_TagsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_measures_or_tags(inputs)
	return __es.explore_no_measures_or_tags(inputs)
});
/**
* | output |
* | --- |
* | "No measures shown" |
*
* @param {Explore_No_Measures_ShownInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_no_measures_shown = /** @type {((inputs?: Explore_No_Measures_ShownInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_No_Measures_ShownInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_no_measures_shown(inputs)
	return __es.explore_no_measures_shown(inputs)
});
/**
* | output |
* | --- |
* | "Only show {tag}" |
*
* @param {Explore_Only_Show_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_only_show_tag = /** @type {((inputs: Explore_Only_Show_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Only_Show_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_only_show_tag(inputs)
	return __es.explore_only_show_tag(inputs)
});
/**
* | output |
* | --- |
* | "Only show this tag" |
*
* @param {Explore_Only_Show_This_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_only_show_this_tag = /** @type {((inputs?: Explore_Only_Show_This_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Only_Show_This_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_only_show_this_tag(inputs)
	return __es.explore_only_show_this_tag(inputs)
});
/**
* | output |
* | --- |
* | "Search dimensions" |
*
* @param {Explore_Search_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_search_dimensions = /** @type {((inputs?: Explore_Search_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Search_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_search_dimensions(inputs)
	return __es.explore_search_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "Search dimensions or tags" |
*
* @param {Explore_Search_Dimensions_Or_TagsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_search_dimensions_or_tags = /** @type {((inputs?: Explore_Search_Dimensions_Or_TagsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Search_Dimensions_Or_TagsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_search_dimensions_or_tags(inputs)
	return __es.explore_search_dimensions_or_tags(inputs)
});
/**
* | output |
* | --- |
* | "Search list" |
*
* @param {Explore_Search_ListInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_search_list = /** @type {((inputs?: Explore_Search_ListInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Search_ListInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_search_list(inputs)
	return __es.explore_search_list(inputs)
});
/**
* | output |
* | --- |
* | "Search measures" |
*
* @param {Explore_Search_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_search_measures = /** @type {((inputs?: Explore_Search_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Search_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_search_measures(inputs)
	return __es.explore_search_measures(inputs)
});
/**
* | output |
* | --- |
* | "Search measures or tags" |
*
* @param {Explore_Search_Measures_Or_TagsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_search_measures_or_tags = /** @type {((inputs?: Explore_Search_Measures_Or_TagsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Search_Measures_Or_TagsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_search_measures_or_tags(inputs)
	return __es.explore_search_measures_or_tags(inputs)
});
/**
* | output |
* | --- |
* | "Show all" |
*
* @param {Explore_Show_AllInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_show_all = /** @type {((inputs?: Explore_Show_AllInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Show_AllInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_show_all(inputs)
	return __es.explore_show_all(inputs)
});
/**
* | output |
* | --- |
* | "Show all in {tag}" |
*
* @param {Explore_Show_All_In_Named_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_show_all_in_named_tag = /** @type {((inputs: Explore_Show_All_In_Named_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Show_All_In_Named_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_show_all_in_named_tag(inputs)
	return __es.explore_show_all_in_named_tag(inputs)
});
/**
* | output |
* | --- |
* | "Show all in tag" |
*
* @param {Explore_Show_All_In_TagInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_show_all_in_tag = /** @type {((inputs?: Explore_Show_All_In_TagInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Show_All_In_TagInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_show_all_in_tag(inputs)
	return __es.explore_show_all_in_tag(inputs)
});
/**
* | output |
* | --- |
* | "Show context for all measures" |
*
* @param {Explore_Show_Context_For_All_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_show_context_for_all_measures = /** @type {((inputs?: Explore_Show_Context_For_All_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Show_Context_For_All_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_show_context_for_all_measures(inputs)
	return __es.explore_show_context_for_all_measures(inputs)
});
/**
* | output |
* | --- |
* | "Show {name}" |
*
* @param {Explore_Show_ItemInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_show_item = /** @type {((inputs: Explore_Show_ItemInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Show_ItemInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_show_item(inputs)
	return __es.explore_show_item(inputs)
});
/**
* | output |
* | --- |
* | "Showing" |
*
* @param {Explore_ShowingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_showing = /** @type {((inputs?: Explore_ShowingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_ShowingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_showing(inputs)
	return __es.explore_showing(inputs)
});
/**
* | output |
* | --- |
* | "Shown dimensions" |
*
* @param {Explore_Shown_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_shown_dimensions = /** @type {((inputs?: Explore_Shown_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Shown_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_shown_dimensions(inputs)
	return __es.explore_shown_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "Shown measures" |
*
* @param {Explore_Shown_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_shown_measures = /** @type {((inputs?: Explore_Shown_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Shown_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_shown_measures(inputs)
	return __es.explore_shown_measures(inputs)
});
/**
* | output |
* | --- |
* | "{visible} of {total} shown" |
*
* @param {Explore_Tag_Shown_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_tag_shown_count = /** @type {((inputs: Explore_Tag_Shown_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Tag_Shown_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_tag_shown_count(inputs)
	return __es.explore_tag_shown_count(inputs)
});
/**
* | output |
* | --- |
* | "Tags" |
*
* @param {Explore_TagsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_tags = /** @type {((inputs?: Explore_TagsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_TagsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_tags(inputs)
	return __es.explore_tags(inputs)
});
/**
* | output |
* | --- |
* | "Unable to open Explore Dashboard" |
*
* @param {Explore_Unable_To_OpenInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_unable_to_open = /** @type {((inputs?: Explore_Unable_To_OpenInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Unable_To_OpenInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_unable_to_open(inputs)
	return __es.explore_unable_to_open(inputs)
});
/**
* | output |
* | --- |
* | "Unknown dimension" |
*
* @param {Explore_Unknown_DimensionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_unknown_dimension = /** @type {((inputs?: Explore_Unknown_DimensionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Unknown_DimensionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_unknown_dimension(inputs)
	return __es.explore_unknown_dimension(inputs)
});
/**
* | output |
* | --- |
* | "Unknown measure" |
*
* @param {Explore_Unknown_MeasureInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const explore_unknown_measure = /** @type {((inputs?: Explore_Unknown_MeasureInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Explore_Unknown_MeasureInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.explore_unknown_measure(inputs)
	return __es.explore_unknown_measure(inputs)
});
/**
* | output |
* | --- |
* | "Add {label} fields" |
*
* @param {Field_List_Add_FieldsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const field_list_add_fields = /** @type {((inputs: Field_List_Add_FieldsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Field_List_Add_FieldsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.field_list_add_fields(inputs)
	return __es.field_list_add_fields(inputs)
});
/**
* | output |
* | --- |
* | "{label} field list" |
*
* @param {Field_List_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const field_list_aria = /** @type {((inputs: Field_List_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Field_List_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.field_list_aria(inputs)
	return __es.field_list_aria(inputs)
});
/**
* | output |
* | --- |
* | "DIMENSIONS" |
*
* @param {Field_List_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const field_list_dimensions = /** @type {((inputs?: Field_List_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Field_List_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.field_list_dimensions(inputs)
	return __es.field_list_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "MEASURES" |
*
* @param {Field_List_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const field_list_measures = /** @type {((inputs?: Field_List_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Field_List_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.field_list_measures(inputs)
	return __es.field_list_measures(inputs)
});
/**
* | output |
* | --- |
* | "TIME" |
*
* @param {Field_List_TimeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const field_list_time = /** @type {((inputs?: Field_List_TimeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Field_List_TimeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.field_list_time(inputs)
	return __es.field_list_time(inputs)
});
/**
* | output |
* | --- |
* | "Advanced (BETA)" |
*
* @param {Filter_Advanced_BetaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_advanced_beta = /** @type {((inputs?: Filter_Advanced_BetaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Advanced_BetaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_advanced_beta(inputs)
	return __es.filter_advanced_beta(inputs)
});
/**
* | output |
* | --- |
* | "Advanced filters are a bleeding edge feature! There may be bugs." |
*
* @param {Filter_Advanced_WarningInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_advanced_warning = /** @type {((inputs?: Filter_Advanced_WarningInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Advanced_WarningInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_advanced_warning(inputs)
	return __es.filter_advanced_warning(inputs)
});
/**
* | output |
* | --- |
* | "DIMENSIONS" |
*
* @param {Filter_DimensionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_dimensions = /** @type {((inputs?: Filter_DimensionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_DimensionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_dimensions(inputs)
	return __es.filter_dimensions(inputs)
});
/**
* | output |
* | --- |
* | "Enter search term" |
*
* @param {Filter_Enter_Search_TermInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_enter_search_term = /** @type {((inputs?: Filter_Enter_Search_TermInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Enter_Search_TermInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_enter_search_term(inputs)
	return __es.filter_enter_search_term(inputs)
});
/**
* | output |
* | --- |
* | "Make filter optional" |
*
* @param {Filter_Make_OptionalInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_make_optional = /** @type {((inputs?: Filter_Make_OptionalInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Make_OptionalInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_make_optional(inputs)
	return __es.filter_make_optional(inputs)
});
/**
* | output |
* | --- |
* | "Make filter required" |
*
* @param {Filter_Make_RequiredInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_make_required = /** @type {((inputs?: Filter_Make_RequiredInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Make_RequiredInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_make_required(inputs)
	return __es.filter_make_required(inputs)
});
/**
* | output |
* | --- |
* | "for {dimension}" |
*
* @param {Filter_Measure_For_DimensionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_for_dimension = /** @type {((inputs: Filter_Measure_For_DimensionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_For_DimensionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_for_dimension(inputs)
	return __es.filter_measure_for_dimension(inputs)
});
/**
* | output |
* | --- |
* | "from {comparison}" |
*
* @param {Filter_Measure_From_ComparisonInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_from_comparison = /** @type {((inputs: Filter_Measure_From_ComparisonInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_From_ComparisonInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_from_comparison(inputs)
	return __es.filter_measure_from_comparison(inputs)
});
/**
* | output |
* | --- |
* | "Between" |
*
* @param {Filter_Measure_Op_BetweenInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_op_between = /** @type {((inputs?: Filter_Measure_Op_BetweenInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_Op_BetweenInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_op_between(inputs)
	return __es.filter_measure_op_between(inputs)
});
/**
* | output |
* | --- |
* | "Does Not Equal" |
*
* @param {Filter_Measure_Op_Does_Not_EqualInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_op_does_not_equal = /** @type {((inputs?: Filter_Measure_Op_Does_Not_EqualInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_Op_Does_Not_EqualInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_op_does_not_equal(inputs)
	return __es.filter_measure_op_does_not_equal(inputs)
});
/**
* | output |
* | --- |
* | "Equals" |
*
* @param {Filter_Measure_Op_EqualsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_op_equals = /** @type {((inputs?: Filter_Measure_Op_EqualsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_Op_EqualsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_op_equals(inputs)
	return __es.filter_measure_op_equals(inputs)
});
/**
* | output |
* | --- |
* | "Greater Than" |
*
* @param {Filter_Measure_Op_Greater_ThanInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_op_greater_than = /** @type {((inputs?: Filter_Measure_Op_Greater_ThanInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_Op_Greater_ThanInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_op_greater_than(inputs)
	return __es.filter_measure_op_greater_than(inputs)
});
/**
* | output |
* | --- |
* | "Greater Than Or Equals" |
*
* @param {Filter_Measure_Op_Greater_Than_Or_EqualsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_op_greater_than_or_equals = /** @type {((inputs?: Filter_Measure_Op_Greater_Than_Or_EqualsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_Op_Greater_Than_Or_EqualsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_op_greater_than_or_equals(inputs)
	return __es.filter_measure_op_greater_than_or_equals(inputs)
});
/**
* | output |
* | --- |
* | "Less Than" |
*
* @param {Filter_Measure_Op_Less_ThanInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_op_less_than = /** @type {((inputs?: Filter_Measure_Op_Less_ThanInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_Op_Less_ThanInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_op_less_than(inputs)
	return __es.filter_measure_op_less_than(inputs)
});
/**
* | output |
* | --- |
* | "Less Than Or Equals" |
*
* @param {Filter_Measure_Op_Less_Than_Or_EqualsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_op_less_than_or_equals = /** @type {((inputs?: Filter_Measure_Op_Less_Than_Or_EqualsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_Op_Less_Than_Or_EqualsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_op_less_than_or_equals(inputs)
	return __es.filter_measure_op_less_than_or_equals(inputs)
});
/**
* | output |
* | --- |
* | "Not Between" |
*
* @param {Filter_Measure_Op_Not_BetweenInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_op_not_between = /** @type {((inputs?: Filter_Measure_Op_Not_BetweenInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_Op_Not_BetweenInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_op_not_between(inputs)
	return __es.filter_measure_op_not_between(inputs)
});
/**
* | output |
* | --- |
* | "change from" |
*
* @param {Filter_Measure_Type_Change_FromInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_type_change_from = /** @type {((inputs?: Filter_Measure_Type_Change_FromInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_Type_Change_FromInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_type_change_from(inputs)
	return __es.filter_measure_type_change_from(inputs)
});
/**
* | output |
* | --- |
* | "% change from" |
*
* @param {Filter_Measure_Type_Percent_Change_FromInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_type_percent_change_from = /** @type {((inputs?: Filter_Measure_Type_Percent_Change_FromInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_Type_Percent_Change_FromInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_type_percent_change_from(inputs)
	return __es.filter_measure_type_percent_change_from(inputs)
});
/**
* | output |
* | --- |
* | "% of total" |
*
* @param {Filter_Measure_Type_Percent_Of_TotalInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_type_percent_of_total = /** @type {((inputs?: Filter_Measure_Type_Percent_Of_TotalInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_Type_Percent_Of_TotalInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_type_percent_of_total(inputs)
	return __es.filter_measure_type_percent_of_total(inputs)
});
/**
* | output |
* | --- |
* | "value" |
*
* @param {Filter_Measure_Type_ValueInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measure_type_value = /** @type {((inputs?: Filter_Measure_Type_ValueInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Measure_Type_ValueInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measure_type_value(inputs)
	return __es.filter_measure_type_value(inputs)
});
/**
* | output |
* | --- |
* | "MEASURES" |
*
* @param {Filter_MeasuresInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_measures = /** @type {((inputs?: Filter_MeasuresInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_MeasuresInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_measures(inputs)
	return __es.filter_measures(inputs)
});
/**
* | output |
* | --- |
* | "Contains" |
*
* @param {Filter_Mode_ContainsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_mode_contains = /** @type {((inputs?: Filter_Mode_ContainsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Mode_ContainsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_mode_contains(inputs)
	return __es.filter_mode_contains(inputs)
});
/**
* | output |
* | --- |
* | "Create a dynamic filter based on a search term" |
*
* @param {Filter_Mode_Contains_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_mode_contains_description = /** @type {((inputs?: Filter_Mode_Contains_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Mode_Contains_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_mode_contains_description(inputs)
	return __es.filter_mode_contains_description(inputs)
});
/**
* | output |
* | --- |
* | "In List" |
*
* @param {Filter_Mode_In_ListInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_mode_in_list = /** @type {((inputs?: Filter_Mode_In_ListInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Mode_In_ListInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_mode_in_list(inputs)
	return __es.filter_mode_in_list(inputs)
});
/**
* | output |
* | --- |
* | "Create a filter based on a list of values" |
*
* @param {Filter_Mode_In_List_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_mode_in_list_description = /** @type {((inputs?: Filter_Mode_In_List_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Mode_In_List_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_mode_in_list_description(inputs)
	return __es.filter_mode_in_list_description(inputs)
});
/**
* | output |
* | --- |
* | "Select" |
*
* @param {Filter_Mode_SelectInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_mode_select = /** @type {((inputs?: Filter_Mode_SelectInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Mode_SelectInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_mode_select(inputs)
	return __es.filter_mode_select(inputs)
});
/**
* | output |
* | --- |
* | "Manually select values for this filter" |
*
* @param {Filter_Mode_Select_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_mode_select_description = /** @type {((inputs?: Filter_Mode_Select_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Mode_Select_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_mode_select_description(inputs)
	return __es.filter_mode_select_description(inputs)
});
/**
* | output |
* | --- |
* | "Paste a list separated by commas or \\n" |
*
* @param {Filter_Paste_List_HintInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_paste_list_hint = /** @type {((inputs?: Filter_Paste_List_HintInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Paste_List_HintInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_paste_list_hint(inputs)
	return __es.filter_paste_list_hint(inputs)
});
/**
* | output |
* | --- |
* | "Pin filter" |
*
* @param {Filter_PinInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_pin = /** @type {((inputs?: Filter_PinInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_PinInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_pin(inputs)
	return __es.filter_pin(inputs)
});
/**
* | output |
* | --- |
* | "Click to pin or unpin : Keep this filter visible at the top so it can't be removed by other users." |
*
* @param {Filter_Pin_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_pin_tooltip = /** @type {((inputs?: Filter_Pin_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Pin_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_pin_tooltip(inputs)
	return __es.filter_pin_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "Click to mark this filter as required. Viewers must set a value for the dashboard to load." |
*
* @param {Filter_Required_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_required_tooltip = /** @type {((inputs?: Filter_Required_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_Required_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_required_tooltip(inputs)
	return __es.filter_required_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "Unpin filter" |
*
* @param {Filter_UnpinInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const filter_unpin = /** @type {((inputs?: Filter_UnpinInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Filter_UnpinInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.filter_unpin(inputs)
	return __es.filter_unpin(inputs)
});
/**
* | output |
* | --- |
* | "Report an issue" |
*
* @param {Footer_Report_IssueInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const footer_report_issue = /** @type {((inputs?: Footer_Report_IssueInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Footer_Report_IssueInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.footer_report_issue(inputs)
	return __es.footer_report_issue(inputs)
});
/**
* | output |
* | --- |
* | "Rill Developer" |
*
* @param {Footer_Rill_DeveloperInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const footer_rill_developer = /** @type {((inputs?: Footer_Rill_DeveloperInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Footer_Rill_DeveloperInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.footer_rill_developer(inputs)
	return __es.footer_rill_developer(inputs)
});
/**
* | output |
* | --- |
* | "Click" |
*
* @param {Footer_Shortcut_ClickInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const footer_shortcut_click = /** @type {((inputs?: Footer_Shortcut_ClickInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Footer_Shortcut_ClickInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.footer_shortcut_click(inputs)
	return __es.footer_shortcut_click(inputs)
});
/**
* | output |
* | --- |
* | "unknown (built from source)" |
*
* @param {Footer_Unknown_VersionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const footer_unknown_version = /** @type {((inputs?: Footer_Unknown_VersionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Footer_Unknown_VersionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.footer_unknown_version(inputs)
	return __es.footer_unknown_version(inputs)
});
/**
* | output |
* | --- |
* | "version" |
*
* @param {Footer_VersionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const footer_version = /** @type {((inputs?: Footer_VersionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Footer_VersionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.footer_version(inputs)
	return __es.footer_version(inputs)
});
/**
* | output |
* | --- |
* | "View documentation" |
*
* @param {Footer_View_DocumentationInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const footer_view_documentation = /** @type {((inputs?: Footer_View_DocumentationInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Footer_View_DocumentationInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.footer_view_documentation(inputs)
	return __es.footer_view_documentation(inputs)
});
/**
* | output |
* | --- |
* | "(optional)" |
*
* @param {Form_OptionalInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const form_optional = /** @type {((inputs?: Form_OptionalInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Form_OptionalInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.form_optional(inputs)
	return __es.form_optional(inputs)
});
/**
* | output |
* | --- |
* | "User group changes saved successfully" |
*
* @param {Groups_Changes_SavedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_changes_saved = /** @type {((inputs?: Groups_Changes_SavedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Changes_SavedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_changes_saved(inputs)
	return __es.groups_changes_saved(inputs)
});
/**
* | output |
* | --- |
* | "Create a group" |
*
* @param {Groups_Create_A_GroupInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_create_a_group = /** @type {((inputs?: Groups_Create_A_GroupInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Create_A_GroupInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_create_a_group(inputs)
	return __es.groups_create_a_group(inputs)
});
/**
* | output |
* | --- |
* | "Create group" |
*
* @param {Groups_Create_GroupInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_create_group = /** @type {((inputs?: Groups_Create_GroupInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Create_GroupInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_create_group(inputs)
	return __es.groups_create_group(inputs)
});
/**
* | output |
* | --- |
* | "User group created" |
*
* @param {Groups_CreatedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_created = /** @type {((inputs?: Groups_CreatedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_CreatedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_created(inputs)
	return __es.groups_created(inputs)
});
/**
* | output |
* | --- |
* | "Delete" |
*
* @param {Groups_DeleteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_delete = /** @type {((inputs?: Groups_DeleteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_DeleteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_delete(inputs)
	return __es.groups_delete(inputs)
});
/**
* | output |
* | --- |
* | "This user group will no longer be able to access the organization." |
*
* @param {Groups_Delete_Confirm_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_delete_confirm_desc = /** @type {((inputs?: Groups_Delete_Confirm_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Delete_Confirm_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_delete_confirm_desc(inputs)
	return __es.groups_delete_confirm_desc(inputs)
});
/**
* | output |
* | --- |
* | "Delete this user group?" |
*
* @param {Groups_Delete_Confirm_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_delete_confirm_title = /** @type {((inputs?: Groups_Delete_Confirm_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Delete_Confirm_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_delete_confirm_title(inputs)
	return __es.groups_delete_confirm_title(inputs)
});
/**
* | output |
* | --- |
* | "User group deleted" |
*
* @param {Groups_DeletedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_deleted = /** @type {((inputs?: Groups_DeletedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_DeletedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_deleted(inputs)
	return __es.groups_deleted(inputs)
});
/**
* | output |
* | --- |
* | "Edit" |
*
* @param {Groups_EditInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_edit = /** @type {((inputs?: Groups_EditInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_EditInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_edit(inputs)
	return __es.groups_edit(inputs)
});
/**
* | output |
* | --- |
* | "Edit group" |
*
* @param {Groups_Edit_GroupInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_edit_group = /** @type {((inputs?: Groups_Edit_GroupInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Edit_GroupInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_edit_group(inputs)
	return __es.groups_edit_group(inputs)
});
/**
* | output |
* | --- |
* | "Error adding role to user group" |
*
* @param {Groups_Error_Adding_RoleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_error_adding_role = /** @type {((inputs?: Groups_Error_Adding_RoleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Error_Adding_RoleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_error_adding_role(inputs)
	return __es.groups_error_adding_role(inputs)
});
/**
* | output |
* | --- |
* | "Error deleting user group" |
*
* @param {Groups_Error_DeletingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_error_deleting = /** @type {((inputs?: Groups_Error_DeletingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Error_DeletingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_error_deleting(inputs)
	return __es.groups_error_deleting(inputs)
});
/**
* | output |
* | --- |
* | "Error loading organization user groups:" |
*
* @param {Groups_Error_LoadingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_error_loading = /** @type {((inputs?: Groups_Error_LoadingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Error_LoadingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_error_loading(inputs)
	return __es.groups_error_loading(inputs)
});
/**
* | output |
* | --- |
* | "Error revoking user group role" |
*
* @param {Groups_Error_Revoking_RoleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_error_revoking_role = /** @type {((inputs?: Groups_Error_Revoking_RoleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Error_Revoking_RoleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_error_revoking_role(inputs)
	return __es.groups_error_revoking_role(inputs)
});
/**
* | output |
* | --- |
* | "Error updating user group role" |
*
* @param {Groups_Error_Updating_RoleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_error_updating_role = /** @type {((inputs?: Groups_Error_Updating_RoleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Error_Updating_RoleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_error_updating_role(inputs)
	return __es.groups_error_updating_role(inputs)
});
/**
* | output |
* | --- |
* | "User group removed" |
*
* @param {Groups_RemovedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_removed = /** @type {((inputs?: Groups_RemovedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_RemovedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_removed(inputs)
	return __es.groups_removed(inputs)
});
/**
* | output |
* | --- |
* | "User group renamed" |
*
* @param {Groups_RenamedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_renamed = /** @type {((inputs?: Groups_RenamedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_RenamedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_renamed(inputs)
	return __es.groups_renamed(inputs)
});
/**
* | output |
* | --- |
* | "User group role added" |
*
* @param {Groups_Role_AddedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_role_added = /** @type {((inputs?: Groups_Role_AddedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Role_AddedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_role_added(inputs)
	return __es.groups_role_added(inputs)
});
/**
* | output |
* | --- |
* | "User group role revoked" |
*
* @param {Groups_Role_RevokedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_role_revoked = /** @type {((inputs?: Groups_Role_RevokedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Role_RevokedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_role_revoked(inputs)
	return __es.groups_role_revoked(inputs)
});
/**
* | output |
* | --- |
* | "User group role updated" |
*
* @param {Groups_Role_UpdatedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_role_updated = /** @type {((inputs?: Groups_Role_UpdatedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Role_UpdatedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_role_updated(inputs)
	return __es.groups_role_updated(inputs)
});
/**
* | output |
* | --- |
* | "No groups found" |
*
* @param {Groups_Table_EmptyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_table_empty = /** @type {((inputs?: Groups_Table_EmptyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Table_EmptyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_table_empty(inputs)
	return __es.groups_table_empty(inputs)
});
/**
* | output |
* | --- |
* | "Group" |
*
* @param {Groups_Table_Header_GroupInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_table_header_group = /** @type {((inputs?: Groups_Table_Header_GroupInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Table_Header_GroupInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_table_header_group(inputs)
	return __es.groups_table_header_group(inputs)
});
/**
* | output |
* | --- |
* | "{count} total groups" |
*
* @param {Groups_Total_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_total_count = /** @type {((inputs: Groups_Total_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Total_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_total_count(inputs)
	return __es.groups_total_count(inputs)
});
/**
* | output |
* | --- |
* | "Yes, delete" |
*
* @param {Groups_Yes_DeleteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const groups_yes_delete = /** @type {((inputs?: Groups_Yes_DeleteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Groups_Yes_DeleteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.groups_yes_delete(inputs)
	return __es.groups_yes_delete(inputs)
});
/**
* | output |
* | --- |
* | "Dashboards" |
*
* @param {Home_Dashboards_HeadingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const home_dashboards_heading = /** @type {((inputs?: Home_Dashboards_HeadingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Home_Dashboards_HeadingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.home_dashboards_heading(inputs)
	return __es.home_dashboards_heading(inputs)
});
/**
* | output |
* | --- |
* | "Explore your dashboards below" |
*
* @param {Home_Subtitle_No_ChatInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const home_subtitle_no_chat = /** @type {((inputs?: Home_Subtitle_No_ChatInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Home_Subtitle_No_ChatInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.home_subtitle_no_chat(inputs)
	return __es.home_subtitle_no_chat(inputs)
});
/**
* | output |
* | --- |
* | "Ask questions about your data, or explore your dashboards below" |
*
* @param {Home_Subtitle_With_ChatInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const home_subtitle_with_chat = /** @type {((inputs?: Home_Subtitle_With_ChatInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Home_Subtitle_With_ChatInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.home_subtitle_with_chat(inputs)
	return __es.home_subtitle_with_chat(inputs)
});
/**
* | output |
* | --- |
* | "Welcome to {projectName}" |
*
* @param {Home_Welcome_ToInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const home_welcome_to = /** @type {((inputs: Home_Welcome_ToInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Home_Welcome_ToInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.home_welcome_to(inputs)
	return __es.home_welcome_to(inputs)
});
/**
* | output |
* | --- |
* | "Day" |
*
* @param {Interval_DayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const interval_day = /** @type {((inputs?: Interval_DayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Interval_DayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.interval_day(inputs)
	return __es.interval_day(inputs)
});
/**
* | output |
* | --- |
* | "Hour" |
*
* @param {Interval_HourInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const interval_hour = /** @type {((inputs?: Interval_HourInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Interval_HourInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.interval_hour(inputs)
	return __es.interval_hour(inputs)
});
/**
* | output |
* | --- |
* | "None" |
*
* @param {Interval_NoneInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const interval_none = /** @type {((inputs?: Interval_NoneInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Interval_NoneInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.interval_none(inputs)
	return __es.interval_none(inputs)
});
/**
* | output |
* | --- |
* | "Week" |
*
* @param {Interval_WeekInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const interval_week = /** @type {((inputs?: Interval_WeekInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Interval_WeekInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.interval_week(inputs)
	return __es.interval_week(inputs)
});
/**
* | output |
* | --- |
* | "no change" |
*
* @param {Kpi_No_ChangeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const kpi_no_change = /** @type {((inputs?: Kpi_No_ChangeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Kpi_No_ChangeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.kpi_no_change(inputs)
	return __es.kpi_no_change(inputs)
});
/**
* | output |
* | --- |
* | "no data" |
*
* @param {Kpi_No_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const kpi_no_data = /** @type {((inputs?: Kpi_No_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Kpi_No_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.kpi_no_data(inputs)
	return __es.kpi_no_data(inputs)
});
/**
* | output |
* | --- |
* | "n/a" |
*
* @param {Kpi_Not_AvailableInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const kpi_not_available = /** @type {((inputs?: Kpi_Not_AvailableInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Kpi_Not_AvailableInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.kpi_not_available(inputs)
	return __es.kpi_not_available(inputs)
});
/**
* | output |
* | --- |
* | "vs {comparison}" |
*
* @param {Kpi_Vs_ComparisonInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const kpi_vs_comparison = /** @type {((inputs: Kpi_Vs_ComparisonInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Kpi_Vs_ComparisonInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.kpi_vs_comparison(inputs)
	return __es.kpi_vs_comparison(inputs)
});
/**
* | output |
* | --- |
* | "English" |
*
* @param {Language_EnInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const language_en = /** @type {((inputs?: Language_EnInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Language_EnInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.language_en(inputs)
	return __es.language_en(inputs)
});
/**
* | output |
* | --- |
* | "Español" |
*
* @param {Language_EsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const language_es = /** @type {((inputs?: Language_EsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Language_EsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.language_es(inputs)
	return __es.language_es(inputs)
});
/**
* | output |
* | --- |
* | "Language" |
*
* @param {Language_Switcher_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const language_switcher_label = /** @type {((inputs?: Language_Switcher_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Language_Switcher_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.language_switcher_label(inputs)
	return __es.language_switcher_label(inputs)
});
/**
* | output |
* | --- |
* | "Inspector Panel" |
*
* @param {Layout_Inspector_Panel_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const layout_inspector_panel_aria = /** @type {((inputs?: Layout_Inspector_Panel_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Layout_Inspector_Panel_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.layout_inspector_panel_aria(inputs)
	return __es.layout_inspector_panel_aria(inputs)
});
/**
* | output |
* | --- |
* | "Compare" |
*
* @param {Leaderboard_CompareInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const leaderboard_compare = /** @type {((inputs?: Leaderboard_CompareInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Leaderboard_CompareInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.leaderboard_compare(inputs)
	return __es.leaderboard_compare(inputs)
});
/**
* | output |
* | --- |
* | "Copy this value to clipboard" |
*
* @param {Leaderboard_Copy_ValueInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const leaderboard_copy_value = /** @type {((inputs?: Leaderboard_Copy_ValueInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Leaderboard_Copy_ValueInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.leaderboard_copy_value(inputs)
	return __es.leaderboard_copy_value(inputs)
});
/**
* | output |
* | --- |
* | "(Expand Table)" |
*
* @param {Leaderboard_Expand_TableInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const leaderboard_expand_table = /** @type {((inputs?: Leaderboard_Expand_TableInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Leaderboard_Expand_TableInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.leaderboard_expand_table(inputs)
	return __es.leaderboard_expand_table(inputs)
});
/**
* | output |
* | --- |
* | "Expand dimension to see more values" |
*
* @param {Leaderboard_Expand_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const leaderboard_expand_tooltip = /** @type {((inputs?: Leaderboard_Expand_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Leaderboard_Expand_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.leaderboard_expand_tooltip(inputs)
	return __es.leaderboard_expand_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "(No available values)" |
*
* @param {Leaderboard_No_Available_ValuesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const leaderboard_no_available_values = /** @type {((inputs?: Leaderboard_No_Available_ValuesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Leaderboard_No_Available_ValuesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.leaderboard_no_available_values(inputs)
	return __es.leaderboard_no_available_values(inputs)
});
/**
* | output |
* | --- |
* | "Remove comparison" |
*
* @param {Leaderboard_Remove_ComparisonInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const leaderboard_remove_comparison = /** @type {((inputs?: Leaderboard_Remove_ComparisonInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Leaderboard_Remove_ComparisonInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.leaderboard_remove_comparison(inputs)
	return __es.leaderboard_remove_comparison(inputs)
});
/**
* | output |
* | --- |
* | "+ Click" |
*
* @param {Leaderboard_Shift_ClickInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const leaderboard_shift_click = /** @type {((inputs?: Leaderboard_Shift_ClickInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Leaderboard_Shift_ClickInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.leaderboard_shift_click(inputs)
	return __es.leaderboard_shift_click(inputs)
});
/**
* | output |
* | --- |
* | "Toggle breakdown for {name} dimension" |
*
* @param {Leaderboard_Toggle_BreakdownInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const leaderboard_toggle_breakdown = /** @type {((inputs: Leaderboard_Toggle_BreakdownInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Leaderboard_Toggle_BreakdownInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.leaderboard_toggle_breakdown(inputs)
	return __es.leaderboard_toggle_breakdown(inputs)
});
/**
* | output |
* | --- |
* | "Add this to your MCP client's configuration file." |
*
* @param {Mcp_Add_To_ConfigInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_add_to_config = /** @type {((inputs?: Mcp_Add_To_ConfigInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Add_To_ConfigInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_add_to_config(inputs)
	return __es.mcp_add_to_config(inputs)
});
/**
* | output |
* | --- |
* | "Add this URL to your AI client's MCP server settings:" |
*
* @param {Mcp_Add_UrlInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_add_url = /** @type {((inputs?: Mcp_Add_UrlInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Add_UrlInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_add_url(inputs)
	return __es.mcp_add_url(inputs)
});
/**
* | output |
* | --- |
* | "Configuration" |
*
* @param {Mcp_ConfigurationInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_configuration = /** @type {((inputs?: Mcp_ConfigurationInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_ConfigurationInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_configuration(inputs)
	return __es.mcp_configuration(inputs)
});
/**
* | output |
* | --- |
* | "Create token" |
*
* @param {Mcp_Create_TokenInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_create_token = /** @type {((inputs?: Mcp_Create_TokenInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Create_TokenInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_create_token(inputs)
	return __es.mcp_create_token(inputs)
});
/**
* | output |
* | --- |
* | "Because this project is {privateLabel}, you need a {tokenLabel} to use in your MCP configuration. This token authenticates your requests." |
*
* @param {Mcp_Create_Token_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_create_token_desc = /** @type {((inputs: Mcp_Create_Token_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Create_Token_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_create_token_desc(inputs)
	return __es.mcp_create_token_desc(inputs)
});
/**
* | output |
* | --- |
* | "Create a personal access token" |
*
* @param {Mcp_Create_Token_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_create_token_title = /** @type {((inputs?: Mcp_Create_Token_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Create_Token_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_create_token_title(inputs)
	return __es.mcp_create_token_title(inputs)
});
/**
* | output |
* | --- |
* | "Ask questions of your Rill project using natural language in any AI client that supports the Model Context Protocol (MCP)." |
*
* @param {Mcp_Dialog_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_dialog_description = /** @type {((inputs?: Mcp_Dialog_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Dialog_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_dialog_description(inputs)
	return __es.mcp_dialog_description(inputs)
});
/**
* | output |
* | --- |
* | "Connect your own AI client" |
*
* @param {Mcp_Dialog_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_dialog_title = /** @type {((inputs?: Mcp_Dialog_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Dialog_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_dialog_title(inputs)
	return __es.mcp_dialog_title(inputs)
});
/**
* | output |
* | --- |
* | "Issuing..." |
*
* @param {Mcp_IssuingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_issuing = /** @type {((inputs?: Mcp_IssuingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_IssuingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_issuing(inputs)
	return __es.mcp_issuing(inputs)
});
/**
* | output |
* | --- |
* | "Learn more" |
*
* @param {Mcp_Learn_MoreInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_learn_more = /** @type {((inputs?: Mcp_Learn_MoreInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Learn_MoreInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_learn_more(inputs)
	return __es.mcp_learn_more(inputs)
});
/**
* | output |
* | --- |
* | "Manual" |
*
* @param {Mcp_Manual_TabInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_manual_tab = /** @type {((inputs?: Mcp_Manual_TabInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Manual_TabInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_manual_tab(inputs)
	return __es.mcp_manual_tab(inputs)
});
/**
* | output |
* | --- |
* | "The OAuth flow will start automatically in your browser." |
*
* @param {Mcp_Oauth_AutoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_oauth_auto = /** @type {((inputs?: Mcp_Oauth_AutoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Oauth_AutoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_oauth_auto(inputs)
	return __es.mcp_oauth_auto(inputs)
});
/**
* | output |
* | --- |
* | "OAuth" |
*
* @param {Mcp_Oauth_TabInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_oauth_tab = /** @type {((inputs?: Mcp_Oauth_TabInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Oauth_TabInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_oauth_tab(inputs)
	return __es.mcp_oauth_tab(inputs)
});
/**
* | output |
* | --- |
* | "personal access token" |
*
* @param {Mcp_Personal_Access_TokenInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_personal_access_token = /** @type {((inputs?: Mcp_Personal_Access_TokenInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Personal_Access_TokenInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_personal_access_token(inputs)
	return __es.mcp_personal_access_token(inputs)
});
/**
* | output |
* | --- |
* | "private" |
*
* @param {Mcp_PrivateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_private = /** @type {((inputs?: Mcp_PrivateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_PrivateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_private(inputs)
	return __es.mcp_private(inputs)
});
/**
* | output |
* | --- |
* | "Recommended" |
*
* @param {Mcp_RecommendedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_recommended = /** @type {((inputs?: Mcp_RecommendedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_RecommendedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_recommended(inputs)
	return __es.mcp_recommended(inputs)
});
/**
* | output |
* | --- |
* | "Token created! Your new token is now included in the configuration snippet below." |
*
* @param {Mcp_Token_CreatedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_token_created = /** @type {((inputs?: Mcp_Token_CreatedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Token_CreatedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_token_created(inputs)
	return __es.mcp_token_created(inputs)
});
/**
* | output |
* | --- |
* | "Failed to issue token. Please try again." |
*
* @param {Mcp_Token_FailedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const mcp_token_failed = /** @type {((inputs?: Mcp_Token_FailedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Mcp_Token_FailedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.mcp_token_failed(inputs)
	return __es.mcp_token_failed(inputs)
});
/**
* | output |
* | --- |
* | "Apply" |
*
* @param {Measure_Filter_ApplyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const measure_filter_apply = /** @type {((inputs?: Measure_Filter_ApplyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Measure_Filter_ApplyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.measure_filter_apply(inputs)
	return __es.measure_filter_apply(inputs)
});
/**
* | output |
* | --- |
* | "By dimension" |
*
* @param {Measure_Filter_By_DimensionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const measure_filter_by_dimension = /** @type {((inputs?: Measure_Filter_By_DimensionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Measure_Filter_By_DimensionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.measure_filter_by_dimension(inputs)
	return __es.measure_filter_by_dimension(inputs)
});
/**
* | output |
* | --- |
* | "Enter a number" |
*
* @param {Measure_Filter_Enter_NumberInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const measure_filter_enter_number = /** @type {((inputs?: Measure_Filter_Enter_NumberInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Measure_Filter_Enter_NumberInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.measure_filter_enter_number(inputs)
	return __es.measure_filter_enter_number(inputs)
});
/**
* | output |
* | --- |
* | "Higher value" |
*
* @param {Measure_Filter_Higher_ValueInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const measure_filter_higher_value = /** @type {((inputs?: Measure_Filter_Higher_ValueInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Measure_Filter_Higher_ValueInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.measure_filter_higher_value(inputs)
	return __es.measure_filter_higher_value(inputs)
});
/**
* | output |
* | --- |
* | "Lower value" |
*
* @param {Measure_Filter_Lower_ValueInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const measure_filter_lower_value = /** @type {((inputs?: Measure_Filter_Lower_ValueInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Measure_Filter_Lower_ValueInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.measure_filter_lower_value(inputs)
	return __es.measure_filter_lower_value(inputs)
});
/**
* | output |
* | --- |
* | "Select dimension to split by" |
*
* @param {Measure_Filter_Select_DimensionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const measure_filter_select_dimension = /** @type {((inputs?: Measure_Filter_Select_DimensionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Measure_Filter_Select_DimensionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.measure_filter_select_dimension(inputs)
	return __es.measure_filter_select_dimension(inputs)
});
/**
* | output |
* | --- |
* | "Threshold" |
*
* @param {Measure_Filter_ThresholdInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const measure_filter_threshold = /** @type {((inputs?: Measure_Filter_ThresholdInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Measure_Filter_ThresholdInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.measure_filter_threshold(inputs)
	return __es.measure_filter_threshold(inputs)
});
/**
* | output |
* | --- |
* | "Choose measures to display" |
*
* @param {Measures_Choose_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const measures_choose_tooltip = /** @type {((inputs?: Measures_Choose_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Measures_Choose_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.measures_choose_tooltip(inputs)
	return __es.measures_choose_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "Measures" |
*
* @param {Measures_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const measures_label = /** @type {((inputs?: Measures_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Measures_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.measures_label(inputs)
	return __es.measures_label(inputs)
});
/**
* | output |
* | --- |
* | "Close sidebar" |
*
* @param {Nav_Close_SidebarInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const nav_close_sidebar = /** @type {((inputs?: Nav_Close_SidebarInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Close_SidebarInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_close_sidebar(inputs)
	return __es.nav_close_sidebar(inputs)
});
/**
* | output |
* | --- |
* | "Data Explorer" |
*
* @param {Nav_Data_ExplorerInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const nav_data_explorer = /** @type {((inputs?: Nav_Data_ExplorerInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Data_ExplorerInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_data_explorer(inputs)
	return __es.nav_data_explorer(inputs)
});
/**
* | output |
* | --- |
* | "Show sidebar" |
*
* @param {Nav_Show_SidebarInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const nav_show_sidebar = /** @type {((inputs?: Nav_Show_SidebarInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Show_SidebarInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_show_sidebar(inputs)
	return __es.nav_show_sidebar(inputs)
});
/**
* | output |
* | --- |
* | "AI" |
*
* @param {Nav_Tab_AiInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const nav_tab_ai = /** @type {((inputs?: Nav_Tab_AiInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Tab_AiInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_tab_ai(inputs)
	return __es.nav_tab_ai(inputs)
});
/**
* | output |
* | --- |
* | "Alerts" |
*
* @param {Nav_Tab_AlertsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const nav_tab_alerts = /** @type {((inputs?: Nav_Tab_AlertsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Tab_AlertsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_tab_alerts(inputs)
	return __es.nav_tab_alerts(inputs)
});
/**
* | output |
* | --- |
* | "Dashboards" |
*
* @param {Nav_Tab_DashboardsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const nav_tab_dashboards = /** @type {((inputs?: Nav_Tab_DashboardsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Tab_DashboardsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_tab_dashboards(inputs)
	return __es.nav_tab_dashboards(inputs)
});
/**
* | output |
* | --- |
* | "Home" |
*
* @param {Nav_Tab_HomeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const nav_tab_home = /** @type {((inputs?: Nav_Tab_HomeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Tab_HomeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_tab_home(inputs)
	return __es.nav_tab_home(inputs)
});
/**
* | output |
* | --- |
* | "Query" |
*
* @param {Nav_Tab_QueryInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const nav_tab_query = /** @type {((inputs?: Nav_Tab_QueryInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Tab_QueryInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_tab_query(inputs)
	return __es.nav_tab_query(inputs)
});
/**
* | output |
* | --- |
* | "Reports" |
*
* @param {Nav_Tab_ReportsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const nav_tab_reports = /** @type {((inputs?: Nav_Tab_ReportsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Tab_ReportsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_tab_reports(inputs)
	return __es.nav_tab_reports(inputs)
});
/**
* | output |
* | --- |
* | "Settings" |
*
* @param {Nav_Tab_SettingsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const nav_tab_settings = /** @type {((inputs?: Nav_Tab_SettingsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Tab_SettingsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_tab_settings(inputs)
	return __es.nav_tab_settings(inputs)
});
/**
* | output |
* | --- |
* | "Status" |
*
* @param {Nav_Tab_StatusInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const nav_tab_status = /** @type {((inputs?: Nav_Tab_StatusInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Nav_Tab_StatusInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.nav_tab_status(inputs)
	return __es.nav_tab_status(inputs)
});
/**
* | output |
* | --- |
* | "Check out your projects below." |
*
* @param {Org_Check_Out_ProjectsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const org_check_out_projects = /** @type {((inputs?: Org_Check_Out_ProjectsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Org_Check_Out_ProjectsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.org_check_out_projects(inputs)
	return __es.org_check_out_projects(inputs)
});
/**
* | output |
* | --- |
* | "+ New project" |
*
* @param {Org_New_ProjectInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const org_new_project = /** @type {((inputs?: Org_New_ProjectInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Org_New_ProjectInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.org_new_project(inputs)
	return __es.org_new_project(inputs)
});
/**
* | output |
* | --- |
* | "Search to add/remove users" |
*
* @param {Org_Search_Add_Remove_UsersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const org_search_add_remove_users = /** @type {((inputs?: Org_Search_Add_Remove_UsersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Org_Search_Add_Remove_UsersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.org_search_add_remove_users(inputs)
	return __es.org_search_add_remove_users(inputs)
});
/**
* | output |
* | --- |
* | "Projects" |
*
* @param {Org_Tab_ProjectsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const org_tab_projects = /** @type {((inputs?: Org_Tab_ProjectsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Org_Tab_ProjectsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.org_tab_projects(inputs)
	return __es.org_tab_projects(inputs)
});
/**
* | output |
* | --- |
* | "Settings" |
*
* @param {Org_Tab_SettingsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const org_tab_settings = /** @type {((inputs?: Org_Tab_SettingsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Org_Tab_SettingsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.org_tab_settings(inputs)
	return __es.org_tab_settings(inputs)
});
/**
* | output |
* | --- |
* | "Users" |
*
* @param {Org_Tab_UsersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const org_tab_users = /** @type {((inputs?: Org_Tab_UsersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Org_Tab_UsersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.org_tab_users(inputs)
	return __es.org_tab_users(inputs)
});
/**
* | output |
* | --- |
* | "Collapse row" |
*
* @param {Pivot_Collapse_RowInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_collapse_row = /** @type {((inputs?: Pivot_Collapse_RowInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Collapse_RowInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_collapse_row(inputs)
	return __es.pivot_collapse_row(inputs)
});
/**
* | output |
* | --- |
* | "dim" |
*
* @param {Pivot_Dim_OneInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_dim_one = /** @type {((inputs?: Pivot_Dim_OneInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Dim_OneInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_dim_one(inputs)
	return __es.pivot_dim_one(inputs)
});
/**
* | output |
* | --- |
* | "dims" |
*
* @param {Pivot_Dim_OtherInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_dim_other = /** @type {((inputs?: Pivot_Dim_OtherInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Dim_OtherInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_dim_other(inputs)
	return __es.pivot_dim_other(inputs)
});
/**
* | output |
* | --- |
* | "Drop here to auto-arrange tag" |
*
* @param {Pivot_Drop_Arrange_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_drop_arrange_aria = /** @type {((inputs?: Pivot_Drop_Arrange_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Drop_Arrange_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_drop_arrange_aria(inputs)
	return __es.pivot_drop_arrange_aria(inputs)
});
/**
* | output |
* | --- |
* | "Drop here to replace rows and columns with this tag" |
*
* @param {Pivot_Drop_Replace_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_drop_replace_aria = /** @type {((inputs?: Pivot_Drop_Replace_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Drop_Replace_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_drop_replace_aria(inputs)
	return __es.pivot_drop_replace_aria(inputs)
});
/**
* | output |
* | --- |
* | "Drop to <strong>replace</strong>: <strong>{dimCount}</strong> {dimLabel} → rows, <strong>{measureCount}</strong> {measureLabel} → columns" |
*
* @param {Pivot_Drop_Replace_TextInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_drop_replace_text = /** @type {((inputs: Pivot_Drop_Replace_TextInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Drop_Replace_TextInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_drop_replace_text(inputs)
	return __es.pivot_drop_replace_text(inputs)
});
/**
* | output |
* | --- |
* | "Drop to replace" |
*
* @param {Pivot_Drop_Split_HintInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_drop_split_hint = /** @type {((inputs?: Pivot_Drop_Split_HintInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Drop_Split_HintInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_drop_split_hint(inputs)
	return __es.pivot_drop_split_hint(inputs)
});
/**
* | output |
* | --- |
* | "Drop here to split: <strong>{dimCount}</strong> {dimLabel} → rows, <strong>{measureCount}</strong> {measureLabel} → columns" |
*
* @param {Pivot_Drop_Split_TextInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_drop_split_text = /** @type {((inputs: Pivot_Drop_Split_TextInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Drop_Split_TextInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_drop_split_text(inputs)
	return __es.pivot_drop_split_text(inputs)
});
/**
* | output |
* | --- |
* | "Expand row" |
*
* @param {Pivot_Expand_RowInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_expand_row = /** @type {((inputs?: Pivot_Expand_RowInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Expand_RowInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_expand_row(inputs)
	return __es.pivot_expand_row(inputs)
});
/**
* | output |
* | --- |
* | "measure" |
*
* @param {Pivot_Measure_OneInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_measure_one = /** @type {((inputs?: Pivot_Measure_OneInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Measure_OneInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_measure_one(inputs)
	return __es.pivot_measure_one(inputs)
});
/**
* | output |
* | --- |
* | "measures" |
*
* @param {Pivot_Measure_OtherInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_measure_other = /** @type {((inputs?: Pivot_Measure_OtherInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Measure_OtherInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_measure_other(inputs)
	return __es.pivot_measure_other(inputs)
});
/**
* | output |
* | --- |
* | "TAGS" |
*
* @param {Pivot_TagsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_tags = /** @type {((inputs?: Pivot_TagsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_TagsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_tags(inputs)
	return __es.pivot_tags(inputs)
});
/**
* | output |
* | --- |
* | "Time {grain}" |
*
* @param {Pivot_Time_Dimension_HeaderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_time_dimension_header = /** @type {((inputs: Pivot_Time_Dimension_HeaderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Time_Dimension_HeaderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_time_dimension_header(inputs)
	return __es.pivot_time_dimension_header(inputs)
});
/**
* | output |
* | --- |
* | "Time" |
*
* @param {Pivot_Time_PrefixInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const pivot_time_prefix = /** @type {((inputs?: Pivot_Time_PrefixInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Pivot_Time_PrefixInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.pivot_time_prefix(inputs)
	return __es.pivot_time_prefix(inputs)
});
/**
* | output |
* | --- |
* | "Project dashboards" |
*
* @param {Project_Dashboards_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const project_dashboards_title = /** @type {((inputs?: Project_Dashboards_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Dashboards_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.project_dashboards_title(inputs)
	return __es.project_dashboards_title(inputs)
});
/**
* | output |
* | --- |
* | "Delete" |
*
* @param {Project_DeleteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const project_delete = /** @type {((inputs?: Project_DeleteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_DeleteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.project_delete(inputs)
	return __es.project_delete(inputs)
});
/**
* | output |
* | --- |
* | "Edit" |
*
* @param {Project_EditInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const project_edit = /** @type {((inputs?: Project_EditInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_EditInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.project_edit(inputs)
	return __es.project_edit(inputs)
});
/**
* | output |
* | --- |
* | "Rename" |
*
* @param {Project_RenameInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const project_rename = /** @type {((inputs?: Project_RenameInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_RenameInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.project_rename(inputs)
	return __es.project_rename(inputs)
});
/**
* | output |
* | --- |
* | "Admin" |
*
* @param {Project_Role_AdminInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const project_role_admin = /** @type {((inputs?: Project_Role_AdminInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Role_AdminInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.project_role_admin(inputs)
	return __es.project_role_admin(inputs)
});
/**
* | output |
* | --- |
* | "Viewer" |
*
* @param {Project_Role_ViewerInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const project_role_viewer = /** @type {((inputs?: Project_Role_ViewerInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Role_ViewerInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.project_role_viewer(inputs)
	return __es.project_role_viewer(inputs)
});
/**
* | output |
* | --- |
* | "Search or invite by email" |
*
* @param {Project_Search_Or_InviteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const project_search_or_invite = /** @type {((inputs?: Project_Search_Or_InviteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Search_Or_InviteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.project_search_or_invite(inputs)
	return __es.project_search_or_invite(inputs)
});
/**
* | output |
* | --- |
* | "Search for users" |
*
* @param {Project_Search_UsersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const project_search_users = /** @type {((inputs?: Project_Search_UsersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Search_UsersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.project_search_users(inputs)
	return __es.project_search_users(inputs)
});
/**
* | output |
* | --- |
* | "Share" |
*
* @param {Project_ShareInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const project_share = /** @type {((inputs?: Project_ShareInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_ShareInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.project_share(inputs)
	return __es.project_share(inputs)
});
/**
* | output |
* | --- |
* | "Share project: {project}" |
*
* @param {Project_Share_HeadingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const project_share_heading = /** @type {((inputs: Project_Share_HeadingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Share_HeadingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.project_share_heading(inputs)
	return __es.project_share_heading(inputs)
});
/**
* | output |
* | --- |
* | "Share project" |
*
* @param {Project_Share_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const project_share_tooltip = /** @type {((inputs?: Project_Share_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Project_Share_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.project_share_tooltip(inputs)
	return __es.project_share_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "Report context menu" |
*
* @param {Report_Context_Menu_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_context_menu_aria = /** @type {((inputs?: Report_Context_Menu_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Context_Menu_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_context_menu_aria(inputs)
	return __es.report_context_menu_aria(inputs)
});
/**
* | output |
* | --- |
* | "Created by {name}" |
*
* @param {Report_Created_ByInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_created_by = /** @type {((inputs: Report_Created_ByInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Created_ByInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_created_by(inputs)
	return __es.report_created_by(inputs)
});
/**
* | output |
* | --- |
* | "Report created through code" |
*
* @param {Report_Created_Through_CodeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_created_through_code = /** @type {((inputs?: Report_Created_Through_CodeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Created_Through_CodeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_created_through_code(inputs)
	return __es.report_created_through_code(inputs)
});
/**
* | output |
* | --- |
* | "Dashboard" |
*
* @param {Report_DashboardInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_dashboard = /** @type {((inputs?: Report_DashboardInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_DashboardInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_dashboard(inputs)
	return __es.report_dashboard(inputs)
});
/**
* | output |
* | --- |
* | "Delete report" |
*
* @param {Report_DeleteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_delete = /** @type {((inputs?: Report_DeleteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_DeleteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_delete(inputs)
	return __es.report_delete(inputs)
});
/**
* | output |
* | --- |
* | "Edit report" |
*
* @param {Report_EditInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_edit = /** @type {((inputs?: Report_EditInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_EditInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_edit(inputs)
	return __es.report_edit(inputs)
});
/**
* | output |
* | --- |
* | "Email recipients" |
*
* @param {Report_Email_RecipientsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_email_recipients = /** @type {((inputs?: Report_Email_RecipientsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Email_RecipientsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_email_recipients(inputs)
	return __es.report_email_recipients(inputs)
});
/**
* | output |
* | --- |
* | "At least one email recipient, slack user, or slack channel is required" |
*
* @param {Report_Email_ValidationInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_email_validation = /** @type {((inputs?: Report_Email_ValidationInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Email_ValidationInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_email_validation(inputs)
	return __es.report_email_validation(inputs)
});
/**
* | output |
* | --- |
* | "Cancel" |
*
* @param {Report_Form_CancelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_cancel = /** @type {((inputs?: Report_Form_CancelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_CancelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_cancel(inputs)
	return __es.report_form_cancel(inputs)
});
/**
* | output |
* | --- |
* | "Channels" |
*
* @param {Report_Form_ChannelsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_channels = /** @type {((inputs?: Report_Form_ChannelsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_ChannelsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_channels(inputs)
	return __es.report_form_channels(inputs)
});
/**
* | output |
* | --- |
* | "Clear filters" |
*
* @param {Report_Form_Clear_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_clear_filters = /** @type {((inputs?: Report_Form_Clear_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Clear_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_clear_filters(inputs)
	return __es.report_form_clear_filters(inputs)
});
/**
* | output |
* | --- |
* | "Columns" |
*
* @param {Report_Form_ColumnsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_columns = /** @type {((inputs?: Report_Form_ColumnsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_ColumnsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_columns(inputs)
	return __es.report_form_columns(inputs)
});
/**
* | output |
* | --- |
* | "Create report" |
*
* @param {Report_Form_CreateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_create = /** @type {((inputs?: Report_Form_CreateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_CreateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_create(inputs)
	return __es.report_form_create(inputs)
});
/**
* | output |
* | --- |
* | "Create" |
*
* @param {Report_Form_Create_ButtonInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_create_button = /** @type {((inputs?: Report_Form_Create_ButtonInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Create_ButtonInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_create_button(inputs)
	return __es.report_form_create_button(inputs)
});
/**
* | output |
* | --- |
* | "Report created" |
*
* @param {Report_Form_Created_NotificationInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_created_notification = /** @type {((inputs?: Report_Form_Created_NotificationInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Created_NotificationInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_created_notification(inputs)
	return __es.report_form_created_notification(inputs)
});
/**
* | output |
* | --- |
* | "Day" |
*
* @param {Report_Form_DayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_day = /** @type {((inputs?: Report_Form_DayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_DayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_day(inputs)
	return __es.report_form_day(inputs)
});
/**
* | output |
* | --- |
* | "First day" |
*
* @param {Report_Form_Day_FirstInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_day_first = /** @type {((inputs?: Report_Form_Day_FirstInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Day_FirstInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_day_first(inputs)
	return __es.report_form_day_first(inputs)
});
/**
* | output |
* | --- |
* | "Friday" |
*
* @param {Report_Form_Day_FridayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_day_friday = /** @type {((inputs?: Report_Form_Day_FridayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Day_FridayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_day_friday(inputs)
	return __es.report_form_day_friday(inputs)
});
/**
* | output |
* | --- |
* | "Monday" |
*
* @param {Report_Form_Day_MondayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_day_monday = /** @type {((inputs?: Report_Form_Day_MondayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Day_MondayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_day_monday(inputs)
	return __es.report_form_day_monday(inputs)
});
/**
* | output |
* | --- |
* | "Saturday" |
*
* @param {Report_Form_Day_SaturdayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_day_saturday = /** @type {((inputs?: Report_Form_Day_SaturdayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Day_SaturdayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_day_saturday(inputs)
	return __es.report_form_day_saturday(inputs)
});
/**
* | output |
* | --- |
* | "Sunday" |
*
* @param {Report_Form_Day_SundayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_day_sunday = /** @type {((inputs?: Report_Form_Day_SundayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Day_SundayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_day_sunday(inputs)
	return __es.report_form_day_sunday(inputs)
});
/**
* | output |
* | --- |
* | "Thursday" |
*
* @param {Report_Form_Day_ThursdayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_day_thursday = /** @type {((inputs?: Report_Form_Day_ThursdayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Day_ThursdayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_day_thursday(inputs)
	return __es.report_form_day_thursday(inputs)
});
/**
* | output |
* | --- |
* | "Tuesday" |
*
* @param {Report_Form_Day_TuesdayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_day_tuesday = /** @type {((inputs?: Report_Form_Day_TuesdayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Day_TuesdayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_day_tuesday(inputs)
	return __es.report_form_day_tuesday(inputs)
});
/**
* | output |
* | --- |
* | "Wednesday" |
*
* @param {Report_Form_Day_WednesdayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_day_wednesday = /** @type {((inputs?: Report_Form_Day_WednesdayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Day_WednesdayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_day_wednesday(inputs)
	return __es.report_form_day_wednesday(inputs)
});
/**
* | output |
* | --- |
* | "docs" |
*
* @param {Report_Form_DocsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_docs = /** @type {((inputs?: Report_Form_DocsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_DocsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_docs(inputs)
	return __es.report_form_docs(inputs)
});
/**
* | output |
* | --- |
* | "Report edited" |
*
* @param {Report_Form_Edited_NotificationInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_edited_notification = /** @type {((inputs?: Report_Form_Edited_NotificationInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Edited_NotificationInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_edited_notification(inputs)
	return __es.report_form_edited_notification(inputs)
});
/**
* | output |
* | --- |
* | "Recipients will receive different views based on their security policy." |
*
* @param {Report_Form_Email_HintInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_email_hint = /** @type {((inputs?: Report_Form_Email_HintInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Email_HintInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_email_hint(inputs)
	return __es.report_form_email_hint(inputs)
});
/**
* | output |
* | --- |
* | "Enter an email address" |
*
* @param {Report_Form_Email_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_email_placeholder = /** @type {((inputs?: Report_Form_Email_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Email_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_email_placeholder(inputs)
	return __es.report_form_email_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Email Recipients" |
*
* @param {Report_Form_Email_RecipientsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_email_recipients = /** @type {((inputs?: Report_Form_Email_RecipientsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Email_RecipientsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_email_recipients(inputs)
	return __es.report_form_email_recipients(inputs)
});
/**
* | output |
* | --- |
* | "Email recurring exports to recipients." |
*
* @param {Report_Form_Email_RecurringInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_email_recurring = /** @type {((inputs?: Report_Form_Email_RecurringInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Email_RecurringInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_email_recurring(inputs)
	return __es.report_form_email_recurring(inputs)
});
/**
* | output |
* | --- |
* | "Filters" |
*
* @param {Report_Form_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_filters = /** @type {((inputs?: Report_Form_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_filters(inputs)
	return __es.report_form_filters(inputs)
});
/**
* | output |
* | --- |
* | "Filters form" |
*
* @param {Report_Form_Filters_AriaInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_filters_aria = /** @type {((inputs?: Report_Form_Filters_AriaInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Filters_AriaInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_filters_aria(inputs)
	return __es.report_form_filters_aria(inputs)
});
/**
* | output |
* | --- |
* | "Format" |
*
* @param {Report_Form_FormatInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_format = /** @type {((inputs?: Report_Form_FormatInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_FormatInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_format(inputs)
	return __es.report_form_format(inputs)
});
/**
* | output |
* | --- |
* | "CSV" |
*
* @param {Report_Form_Format_CsvInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_format_csv = /** @type {((inputs?: Report_Form_Format_CsvInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Format_CsvInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_format_csv(inputs)
	return __es.report_form_format_csv(inputs)
});
/**
* | output |
* | --- |
* | "Parquet" |
*
* @param {Report_Form_Format_ParquetInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_format_parquet = /** @type {((inputs?: Report_Form_Format_ParquetInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Format_ParquetInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_format_parquet(inputs)
	return __es.report_form_format_parquet(inputs)
});
/**
* | output |
* | --- |
* | "XLSX" |
*
* @param {Report_Form_Format_XlsxInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_format_xlsx = /** @type {((inputs?: Report_Form_Format_XlsxInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Format_XlsxInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_format_xlsx(inputs)
	return __es.report_form_format_xlsx(inputs)
});
/**
* | output |
* | --- |
* | "Daily" |
*
* @param {Report_Form_Freq_DailyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_freq_daily = /** @type {((inputs?: Report_Form_Freq_DailyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Freq_DailyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_freq_daily(inputs)
	return __es.report_form_freq_daily(inputs)
});
/**
* | output |
* | --- |
* | "Monthly" |
*
* @param {Report_Form_Freq_MonthlyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_freq_monthly = /** @type {((inputs?: Report_Form_Freq_MonthlyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Freq_MonthlyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_freq_monthly(inputs)
	return __es.report_form_freq_monthly(inputs)
});
/**
* | output |
* | --- |
* | "Weekdays" |
*
* @param {Report_Form_Freq_WeekdaysInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_freq_weekdays = /** @type {((inputs?: Report_Form_Freq_WeekdaysInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Freq_WeekdaysInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_freq_weekdays(inputs)
	return __es.report_form_freq_weekdays(inputs)
});
/**
* | output |
* | --- |
* | "Weekly" |
*
* @param {Report_Form_Freq_WeeklyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_freq_weekly = /** @type {((inputs?: Report_Form_Freq_WeeklyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Freq_WeeklyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_freq_weekly(inputs)
	return __es.report_form_freq_weekly(inputs)
});
/**
* | output |
* | --- |
* | "Frequency" |
*
* @param {Report_Form_FrequencyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_frequency = /** @type {((inputs?: Report_Form_FrequencyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_FrequencyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_frequency(inputs)
	return __es.report_form_frequency(inputs)
});
/**
* | output |
* | --- |
* | "Go to scheduled reports" |
*
* @param {Report_Form_Go_To_ReportsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_go_to_reports = /** @type {((inputs?: Report_Form_Go_To_ReportsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Go_To_ReportsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_go_to_reports(inputs)
	return __es.report_form_go_to_reports(inputs)
});
/**
* | output |
* | --- |
* | "Include metadata" |
*
* @param {Report_Form_Include_MetadataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_include_metadata = /** @type {((inputs?: Report_Form_Include_MetadataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Include_MetadataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_include_metadata(inputs)
	return __es.report_form_include_metadata(inputs)
});
/**
* | output |
* | --- |
* | "Invalid email" |
*
* @param {Report_Form_Invalid_EmailInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_invalid_email = /** @type {((inputs?: Report_Form_Invalid_EmailInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Invalid_EmailInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_invalid_email(inputs)
	return __es.report_form_invalid_email(inputs)
});
/**
* | output |
* | --- |
* | "(Local)" |
*
* @param {Report_Form_Local_SuffixInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_local_suffix = /** @type {((inputs?: Report_Form_Local_SuffixInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Local_SuffixInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_local_suffix(inputs)
	return __es.report_form_local_suffix(inputs)
});
/**
* | output |
* | --- |
* | "Adds a header to the file that includes filters, time range, and other metadata." |
*
* @param {Report_Form_Metadata_TooltipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_metadata_tooltip = /** @type {((inputs?: Report_Form_Metadata_TooltipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Metadata_TooltipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_metadata_tooltip(inputs)
	return __es.report_form_metadata_tooltip(inputs)
});
/**
* | output |
* | --- |
* | "No rows selected" |
*
* @param {Report_Form_No_Rows_SelectedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_no_rows_selected = /** @type {((inputs?: Report_Form_No_Rows_SelectedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_No_Rows_SelectedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_no_rows_selected(inputs)
	return __es.report_form_no_rows_selected(inputs)
});
/**
* | output |
* | --- |
* | "Recipients must be part of the project when running as recipient" |
*
* @param {Report_Form_Recipients_Must_Be_ProjectInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_recipients_must_be_project = /** @type {((inputs?: Report_Form_Recipients_Must_Be_ProjectInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Recipients_Must_Be_ProjectInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_recipients_must_be_project(inputs)
	return __es.report_form_recipients_must_be_project(inputs)
});
/**
* | output |
* | --- |
* | "Required" |
*
* @param {Report_Form_RequiredInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_required = /** @type {((inputs?: Report_Form_RequiredInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_RequiredInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_required(inputs)
	return __es.report_form_required(inputs)
});
/**
* | output |
* | --- |
* | "Row limit" |
*
* @param {Report_Form_Row_LimitInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_row_limit = /** @type {((inputs?: Report_Form_Row_LimitInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Row_LimitInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_row_limit(inputs)
	return __es.report_form_row_limit(inputs)
});
/**
* | output |
* | --- |
* | "1000" |
*
* @param {Report_Form_Row_Limit_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_row_limit_placeholder = /** @type {((inputs?: Report_Form_Row_Limit_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Row_Limit_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_row_limit_placeholder(inputs)
	return __es.report_form_row_limit_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Rows" |
*
* @param {Report_Form_RowsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_rows = /** @type {((inputs?: Report_Form_RowsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_RowsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_rows(inputs)
	return __es.report_form_rows(inputs)
});
/**
* | output |
* | --- |
* | "Run as" |
*
* @param {Report_Form_Run_AsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_run_as = /** @type {((inputs?: Report_Form_Run_AsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Run_AsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_run_as(inputs)
	return __es.report_form_run_as(inputs)
});
/**
* | output |
* | --- |
* | "Creator" |
*
* @param {Report_Form_Run_As_CreatorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_run_as_creator = /** @type {((inputs?: Report_Form_Run_As_CreatorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Run_As_CreatorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_run_as_creator(inputs)
	return __es.report_form_run_as_creator(inputs)
});
/**
* | output |
* | --- |
* | "Works for any recipient, including external recipient. It does NOT grant access beyond the report's filters and dashboard." |
*
* @param {Report_Form_Run_As_Creator_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_run_as_creator_desc = /** @type {((inputs?: Report_Form_Run_As_Creator_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Run_As_Creator_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_run_as_creator_desc(inputs)
	return __es.report_form_run_as_creator_desc(inputs)
});
/**
* | output |
* | --- |
* | "Recipient" |
*
* @param {Report_Form_Run_As_RecipientInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_run_as_recipient = /** @type {((inputs?: Report_Form_Run_As_RecipientInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Run_As_RecipientInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_run_as_recipient(inputs)
	return __es.report_form_run_as_recipient(inputs)
});
/**
* | output |
* | --- |
* | "Does NOT work for non-project members." |
*
* @param {Report_Form_Run_As_Recipient_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_run_as_recipient_desc = /** @type {((inputs?: Report_Form_Run_As_Recipient_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Run_As_Recipient_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_run_as_recipient_desc(inputs)
	return __es.report_form_run_as_recipient_desc(inputs)
});
/**
* | output |
* | --- |
* | "Save report" |
*
* @param {Report_Form_SaveInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_save = /** @type {((inputs?: Report_Form_SaveInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_SaveInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_save(inputs)
	return __es.report_form_save(inputs)
});
/**
* | output |
* | --- |
* | "Save" |
*
* @param {Report_Form_Save_ButtonInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_save_button = /** @type {((inputs?: Report_Form_Save_ButtonInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Save_ButtonInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_save_button(inputs)
	return __es.report_form_save_button(inputs)
});
/**
* | output |
* | --- |
* | "Schedule report" |
*
* @param {Report_Form_ScheduleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_schedule = /** @type {((inputs?: Report_Form_ScheduleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_ScheduleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_schedule(inputs)
	return __es.report_form_schedule(inputs)
});
/**
* | output |
* | --- |
* | "We'll send alerts directly to these channels." |
*
* @param {Report_Form_Slack_Channels_HintInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_slack_channels_hint = /** @type {((inputs?: Report_Form_Slack_Channels_HintInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Slack_Channels_HintInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_slack_channels_hint(inputs)
	return __es.report_form_slack_channels_hint(inputs)
});
/**
* | output |
* | --- |
* | "Slack has not been configured for this project. Read the {link} to learn more." |
*
* @param {Report_Form_Slack_Not_ConfiguredInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_slack_not_configured = /** @type {((inputs: Report_Form_Slack_Not_ConfiguredInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Slack_Not_ConfiguredInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_slack_not_configured(inputs)
	return __es.report_form_slack_not_configured(inputs)
});
/**
* | output |
* | --- |
* | "Slack notifications" |
*
* @param {Report_Form_Slack_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_slack_title = /** @type {((inputs?: Report_Form_Slack_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Slack_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_slack_title(inputs)
	return __es.report_form_slack_title(inputs)
});
/**
* | output |
* | --- |
* | "Users" |
*
* @param {Report_Form_Slack_UsersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_slack_users = /** @type {((inputs?: Report_Form_Slack_UsersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Slack_UsersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_slack_users(inputs)
	return __es.report_form_slack_users(inputs)
});
/**
* | output |
* | --- |
* | "We'll alert them with direct messages in Slack." |
*
* @param {Report_Form_Slack_Users_HintInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_slack_users_hint = /** @type {((inputs?: Report_Form_Slack_Users_HintInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Slack_Users_HintInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_slack_users_hint(inputs)
	return __es.report_form_slack_users_hint(inputs)
});
/**
* | output |
* | --- |
* | "Time" |
*
* @param {Report_Form_TimeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_time = /** @type {((inputs?: Report_Form_TimeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_TimeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_time(inputs)
	return __es.report_form_time(inputs)
});
/**
* | output |
* | --- |
* | "Time zone" |
*
* @param {Report_Form_TimezoneInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_timezone = /** @type {((inputs?: Report_Form_TimezoneInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_TimezoneInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_timezone(inputs)
	return __es.report_form_timezone(inputs)
});
/**
* | output |
* | --- |
* | "Report title" |
*
* @param {Report_Form_Title_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_title_label = /** @type {((inputs?: Report_Form_Title_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Title_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_title_label(inputs)
	return __es.report_form_title_label(inputs)
});
/**
* | output |
* | --- |
* | "My report" |
*
* @param {Report_Form_Title_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_form_title_placeholder = /** @type {((inputs?: Report_Form_Title_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Form_Title_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_form_title_placeholder(inputs)
	return __es.report_form_title_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Name" |
*
* @param {Report_Name_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_name_label = /** @type {((inputs?: Report_Name_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Name_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_name_label(inputs)
	return __es.report_name_label(inputs)
});
/**
* | output |
* | --- |
* | "Next run" |
*
* @param {Report_Next_RunInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_next_run = /** @type {((inputs?: Report_Next_RunInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Next_RunInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_next_run(inputs)
	return __es.report_next_run(inputs)
});
/**
* | output |
* | --- |
* | "No row limit" |
*
* @param {Report_No_Row_LimitInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_no_row_limit = /** @type {((inputs?: Report_No_Row_LimitInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_No_Row_LimitInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_no_row_limit(inputs)
	return __es.report_no_row_limit(inputs)
});
/**
* | output |
* | --- |
* | "Repeats" |
*
* @param {Report_RepeatsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_repeats = /** @type {((inputs?: Report_RepeatsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_RepeatsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_repeats(inputs)
	return __es.report_repeats(inputs)
});
/**
* | output |
* | --- |
* | "Slack recipients" |
*
* @param {Report_Slack_RecipientsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_slack_recipients = /** @type {((inputs?: Report_Slack_RecipientsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Slack_RecipientsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_slack_recipients(inputs)
	return __es.report_slack_recipients(inputs)
});
/**
* | output |
* | --- |
* | "Failed" |
*
* @param {Report_Status_FailedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_status_failed = /** @type {((inputs?: Report_Status_FailedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Status_FailedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_status_failed(inputs)
	return __es.report_status_failed(inputs)
});
/**
* | output |
* | --- |
* | "Report sent" |
*
* @param {Report_Status_SentInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_status_sent = /** @type {((inputs?: Report_Status_SentInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Status_SentInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_status_sent(inputs)
	return __es.report_status_sent(inputs)
});
/**
* | output |
* | --- |
* | "Triggered an ad-hoc run of this report." |
*
* @param {Report_Triggered_AdhocInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_triggered_adhoc = /** @type {((inputs?: Report_Triggered_AdhocInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Triggered_AdhocInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_triggered_adhoc(inputs)
	return __es.report_triggered_adhoc(inputs)
});
/**
* | output |
* | --- |
* | "Failed to unsubscribe." |
*
* @param {Report_Unsubscribe_FailedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_unsubscribe_failed = /** @type {((inputs?: Report_Unsubscribe_FailedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_Unsubscribe_FailedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_unsubscribe_failed(inputs)
	return __es.report_unsubscribe_failed(inputs)
});
/**
* | output |
* | --- |
* | "Unsubscribed from report." |
*
* @param {Report_UnsubscribedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_unsubscribed = /** @type {((inputs?: Report_UnsubscribedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_UnsubscribedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_unsubscribed(inputs)
	return __es.report_unsubscribed(inputs)
});
/**
* | output |
* | --- |
* | "Unsubscribing..." |
*
* @param {Report_UnsubscribingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const report_unsubscribing = /** @type {((inputs?: Report_UnsubscribingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Report_UnsubscribingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.report_unsubscribing(inputs)
	return __es.report_unsubscribing(inputs)
});
/**
* | output |
* | --- |
* | "Schedule {reportsLink} from any dashboard" |
*
* @param {Reports_Empty_ActionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const reports_empty_action = /** @type {((inputs: Reports_Empty_ActionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Reports_Empty_ActionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.reports_empty_action(inputs)
	return __es.reports_empty_action(inputs)
});
/**
* | output |
* | --- |
* | "You don't have any reports yet" |
*
* @param {Reports_Empty_MessageInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const reports_empty_message = /** @type {((inputs?: Reports_Empty_MessageInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Reports_Empty_MessageInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.reports_empty_message(inputs)
	return __es.reports_empty_message(inputs)
});
/**
* | output |
* | --- |
* | "reports" |
*
* @param {Reports_Link_TextInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const reports_link_text = /** @type {((inputs?: Reports_Link_TextInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Reports_Link_TextInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.reports_link_text(inputs)
	return __es.reports_link_text(inputs)
});
/**
* | output |
* | --- |
* | "If this error persists, please contact support." |
*
* @param {Resource_Error_Contact_SupportInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const resource_error_contact_support = /** @type {((inputs?: Resource_Error_Contact_SupportInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Resource_Error_Contact_SupportInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.resource_error_contact_support(inputs)
	return __es.resource_error_contact_support(inputs)
});
/**
* | output |
* | --- |
* | "Error loading {kind}s" |
*
* @param {Resource_Error_LoadingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const resource_error_loading = /** @type {((inputs: Resource_Error_LoadingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Resource_Error_LoadingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.resource_error_loading(inputs)
	return __es.resource_error_loading(inputs)
});
/**
* | output |
* | --- |
* | "Search" |
*
* @param {Resource_Search_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const resource_search_placeholder = /** @type {((inputs?: Resource_Search_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Resource_Search_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.resource_search_placeholder(inputs)
	return __es.resource_search_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Canvas" |
*
* @param {Resource_Type_CanvasInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const resource_type_canvas = /** @type {((inputs?: Resource_Type_CanvasInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Resource_Type_CanvasInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.resource_type_canvas(inputs)
	return __es.resource_type_canvas(inputs)
});
/**
* | output |
* | --- |
* | "Explore" |
*
* @param {Resource_Type_ExploreInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const resource_type_explore = /** @type {((inputs?: Resource_Type_ExploreInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Resource_Type_ExploreInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.resource_type_explore(inputs)
	return __es.resource_type_explore(inputs)
});
/**
* | output |
* | --- |
* | "Admin" |
*
* @param {Role_AdminInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_admin = /** @type {((inputs?: Role_AdminInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_AdminInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_admin(inputs)
	return __es.role_admin(inputs)
});
/**
* | output |
* | --- |
* | "Editor" |
*
* @param {Role_EditorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_editor = /** @type {((inputs?: Role_EditorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_EditorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_editor(inputs)
	return __es.role_editor(inputs)
});
/**
* | output |
* | --- |
* | "Guest" |
*
* @param {Role_GuestInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_guest = /** @type {((inputs?: Role_GuestInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_GuestInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_guest(inputs)
	return __es.role_guest(inputs)
});
/**
* | output |
* | --- |
* | "Access to invited projects only" |
*
* @param {Role_Guest_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_guest_desc = /** @type {((inputs?: Role_Guest_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_Guest_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_guest_desc(inputs)
	return __es.role_guest_desc(inputs)
});
/**
* | output |
* | --- |
* | "Full control over organization settings and members" |
*
* @param {Role_Org_Admin_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_org_admin_desc = /** @type {((inputs?: Role_Org_Admin_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_Org_Admin_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_org_admin_desc(inputs)
	return __es.role_org_admin_desc(inputs)
});
/**
* | output |
* | --- |
* | "Full access to org settings, members, and all projects" |
*
* @param {Role_Org_Admin_DetailInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_org_admin_detail = /** @type {((inputs?: Role_Org_Admin_DetailInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_Org_Admin_DetailInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_org_admin_detail(inputs)
	return __es.role_org_admin_detail(inputs)
});
/**
* | output |
* | --- |
* | "Can manage projects and most org resources" |
*
* @param {Role_Org_Editor_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_org_editor_desc = /** @type {((inputs?: Role_Org_Editor_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_Org_Editor_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_org_editor_desc(inputs)
	return __es.role_org_editor_desc(inputs)
});
/**
* | output |
* | --- |
* | "Can create/manage projects and non-admin members" |
*
* @param {Role_Org_Editor_DetailInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_org_editor_detail = /** @type {((inputs?: Role_Org_Editor_DetailInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_Org_Editor_DetailInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_org_editor_detail(inputs)
	return __es.role_org_editor_detail(inputs)
});
/**
* | output |
* | --- |
* | "Read-only access to organization and projects" |
*
* @param {Role_Org_Viewer_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_org_viewer_desc = /** @type {((inputs?: Role_Org_Viewer_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_Org_Viewer_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_org_viewer_desc(inputs)
	return __es.role_org_viewer_desc(inputs)
});
/**
* | output |
* | --- |
* | "Read-only access to all org projects" |
*
* @param {Role_Org_Viewer_DetailInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_org_viewer_detail = /** @type {((inputs?: Role_Org_Viewer_DetailInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_Org_Viewer_DetailInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_org_viewer_detail(inputs)
	return __es.role_org_viewer_detail(inputs)
});
/**
* | output |
* | --- |
* | "Full control of project settings and members" |
*
* @param {Role_Project_Admin_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_project_admin_desc = /** @type {((inputs?: Role_Project_Admin_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_Project_Admin_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_project_admin_desc(inputs)
	return __es.role_project_admin_desc(inputs)
});
/**
* | output |
* | --- |
* | "Can create and edit dashboards; manage non-admin access" |
*
* @param {Role_Project_Editor_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_project_editor_desc = /** @type {((inputs?: Role_Project_Editor_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_Project_Editor_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_project_editor_desc(inputs)
	return __es.role_project_editor_desc(inputs)
});
/**
* | output |
* | --- |
* | "Read-only access to all project resources" |
*
* @param {Role_Project_Viewer_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_project_viewer_desc = /** @type {((inputs?: Role_Project_Viewer_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_Project_Viewer_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_project_viewer_desc(inputs)
	return __es.role_project_viewer_desc(inputs)
});
/**
* | output |
* | --- |
* | "Viewer" |
*
* @param {Role_ViewerInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const role_viewer = /** @type {((inputs?: Role_ViewerInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Role_ViewerInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.role_viewer(inputs)
	return __es.role_viewer(inputs)
});
/**
* | output |
* | --- |
* | "Limited view. For full access and features, visit the {link}." |
*
* @param {Share_Limited_ViewInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const share_limited_view = /** @type {((inputs: Share_Limited_ViewInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Share_Limited_ViewInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.share_limited_view(inputs)
	return __es.share_limited_view(inputs)
});
/**
* | output |
* | --- |
* | "original dashboard" |
*
* @param {Share_Original_DashboardInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const share_original_dashboard = /** @type {((inputs?: Share_Original_DashboardInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Share_Original_DashboardInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.share_original_dashboard(inputs)
	return __es.share_original_dashboard(inputs)
});
/**
* | output |
* | --- |
* | "Off" |
*
* @param {Snooze_OffInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const snooze_off = /** @type {((inputs?: Snooze_OffInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Snooze_OffInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.snooze_off(inputs)
	return __es.snooze_off(inputs)
});
/**
* | output |
* | --- |
* | "Describe" |
*
* @param {Status_Action_DescribeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_action_describe = /** @type {((inputs?: Status_Action_DescribeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Action_DescribeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_action_describe(inputs)
	return __es.status_action_describe(inputs)
});
/**
* | output |
* | --- |
* | "Full Refresh" |
*
* @param {Status_Action_Full_RefreshInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_action_full_refresh = /** @type {((inputs?: Status_Action_Full_RefreshInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Action_Full_RefreshInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_action_full_refresh(inputs)
	return __es.status_action_full_refresh(inputs)
});
/**
* | output |
* | --- |
* | "Incremental Refresh" |
*
* @param {Status_Action_Incremental_RefreshInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_action_incremental_refresh = /** @type {((inputs?: Status_Action_Incremental_RefreshInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Action_Incremental_RefreshInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_action_incremental_refresh(inputs)
	return __es.status_action_incremental_refresh(inputs)
});
/**
* | output |
* | --- |
* | "Refresh Errored Partitions" |
*
* @param {Status_Action_Refresh_Errored_PartitionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_action_refresh_errored_partitions = /** @type {((inputs?: Status_Action_Refresh_Errored_PartitionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Action_Refresh_Errored_PartitionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_action_refresh_errored_partitions(inputs)
	return __es.status_action_refresh_errored_partitions(inputs)
});
/**
* | output |
* | --- |
* | "View Logs" |
*
* @param {Status_Action_View_LogsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_action_view_logs = /** @type {((inputs?: Status_Action_View_LogsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Action_View_LogsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_action_view_logs(inputs)
	return __es.status_action_view_logs(inputs)
});
/**
* | output |
* | --- |
* | "View Partitions" |
*
* @param {Status_Action_View_PartitionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_action_view_partitions = /** @type {((inputs?: Status_Action_View_PartitionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Action_View_PartitionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_action_view_partitions(inputs)
	return __es.status_action_view_partitions(inputs)
});
/**
* | output |
* | --- |
* | "Cancel" |
*
* @param {Status_CancelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_cancel = /** @type {((inputs?: Status_CancelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_CancelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_cancel(inputs)
	return __es.status_cancel(inputs)
});
/**
* | output |
* | --- |
* | "Checking for errors..." |
*
* @param {Status_Checking_ErrorsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_checking_errors = /** @type {((inputs?: Status_Checking_ErrorsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Checking_ErrorsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_checking_errors(inputs)
	return __es.status_checking_errors(inputs)
});
/**
* | output |
* | --- |
* | "Clear" |
*
* @param {Status_ClearInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_clear = /** @type {((inputs?: Status_ClearInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_ClearInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_clear(inputs)
	return __es.status_clear(inputs)
});
/**
* | output |
* | --- |
* | "Clone this project to develop locally." |
*
* @param {Status_Clone_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_clone_description = /** @type {((inputs?: Status_Clone_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Clone_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_clone_description(inputs)
	return __es.status_clone_description(inputs)
});
/**
* | output |
* | --- |
* | "Database Size" |
*
* @param {Status_Column_Database_SizeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_column_database_size = /** @type {((inputs?: Status_Column_Database_SizeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Column_Database_SizeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_column_database_size(inputs)
	return __es.status_column_database_size(inputs)
});
/**
* | output |
* | --- |
* | "Last refresh" |
*
* @param {Status_Column_Last_RefreshInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_column_last_refresh = /** @type {((inputs?: Status_Column_Last_RefreshInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Column_Last_RefreshInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_column_last_refresh(inputs)
	return __es.status_column_last_refresh(inputs)
});
/**
* | output |
* | --- |
* | "Model Name" |
*
* @param {Status_Column_Model_NameInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_column_model_name = /** @type {((inputs?: Status_Column_Model_NameInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Column_Model_NameInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_column_model_name(inputs)
	return __es.status_column_model_name(inputs)
});
/**
* | output |
* | --- |
* | "Name" |
*
* @param {Status_Column_NameInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_column_name = /** @type {((inputs?: Status_Column_NameInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Column_NameInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_column_name(inputs)
	return __es.status_column_name(inputs)
});
/**
* | output |
* | --- |
* | "Next refresh" |
*
* @param {Status_Column_Next_RefreshInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_column_next_refresh = /** @type {((inputs?: Status_Column_Next_RefreshInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Column_Next_RefreshInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_column_next_refresh(inputs)
	return __es.status_column_next_refresh(inputs)
});
/**
* | output |
* | --- |
* | "Table Name" |
*
* @param {Status_Column_Table_NameInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_column_table_name = /** @type {((inputs?: Status_Column_Table_NameInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Column_Table_NameInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_column_table_name(inputs)
	return __es.status_column_table_name(inputs)
});
/**
* | output |
* | --- |
* | "Type" |
*
* @param {Status_Column_TypeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_column_type = /** @type {((inputs?: Status_Column_TypeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Column_TypeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_column_type(inputs)
	return __es.status_column_type(inputs)
});
/**
* | output |
* | --- |
* | "Complete" |
*
* @param {Status_CompleteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_complete = /** @type {((inputs?: Status_CompleteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_CompleteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_complete(inputs)
	return __es.status_complete(inputs)
});
/**
* | output |
* | --- |
* | "Compute unit" |
*
* @param {Status_Compute_UnitInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_compute_unit = /** @type {((inputs?: Status_Compute_UnitInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Compute_UnitInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_compute_unit(inputs)
	return __es.status_compute_unit(inputs)
});
/**
* | output |
* | --- |
* | "Compute units" |
*
* @param {Status_Compute_UnitsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_compute_units = /** @type {((inputs?: Status_Compute_UnitsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Compute_UnitsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_compute_units(inputs)
	return __es.status_compute_units(inputs)
});
/**
* | output |
* | --- |
* | "Data accessible" |
*
* @param {Status_Data_AccessibleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_data_accessible = /** @type {((inputs?: Status_Data_AccessibleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Data_AccessibleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_data_accessible(inputs)
	return __es.status_data_accessible(inputs)
});
/**
* | output |
* | --- |
* | "Data size" |
*
* @param {Status_Data_SizeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_data_size = /** @type {((inputs?: Status_Data_SizeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Data_SizeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_data_size(inputs)
	return __es.status_data_size(inputs)
});
/**
* | output |
* | --- |
* | "Deleted" |
*
* @param {Status_Deploy_DeletedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_deploy_deleted = /** @type {((inputs?: Status_Deploy_DeletedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Deploy_DeletedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_deploy_deleted(inputs)
	return __es.status_deploy_deleted(inputs)
});
/**
* | output |
* | --- |
* | "Deleting" |
*
* @param {Status_Deploy_DeletingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_deploy_deleting = /** @type {((inputs?: Status_Deploy_DeletingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Deploy_DeletingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_deploy_deleting(inputs)
	return __es.status_deploy_deleting(inputs)
});
/**
* | output |
* | --- |
* | "Error" |
*
* @param {Status_Deploy_ErrorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_deploy_error = /** @type {((inputs?: Status_Deploy_ErrorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Deploy_ErrorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_deploy_error(inputs)
	return __es.status_deploy_error(inputs)
});
/**
* | output |
* | --- |
* | "Not deployed" |
*
* @param {Status_Deploy_Not_DeployedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_deploy_not_deployed = /** @type {((inputs?: Status_Deploy_Not_DeployedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Deploy_Not_DeployedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_deploy_not_deployed(inputs)
	return __es.status_deploy_not_deployed(inputs)
});
/**
* | output |
* | --- |
* | "Pending" |
*
* @param {Status_Deploy_PendingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_deploy_pending = /** @type {((inputs?: Status_Deploy_PendingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Deploy_PendingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_deploy_pending(inputs)
	return __es.status_deploy_pending(inputs)
});
/**
* | output |
* | --- |
* | "Ready" |
*
* @param {Status_Deploy_ReadyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_deploy_ready = /** @type {((inputs?: Status_Deploy_ReadyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Deploy_ReadyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_deploy_ready(inputs)
	return __es.status_deploy_ready(inputs)
});
/**
* | output |
* | --- |
* | "Stopped" |
*
* @param {Status_Deploy_StoppedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_deploy_stopped = /** @type {((inputs?: Status_Deploy_StoppedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Deploy_StoppedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_deploy_stopped(inputs)
	return __es.status_deploy_stopped(inputs)
});
/**
* | output |
* | --- |
* | "Stopping" |
*
* @param {Status_Deploy_StoppingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_deploy_stopping = /** @type {((inputs?: Status_Deploy_StoppingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Deploy_StoppingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_deploy_stopping(inputs)
	return __es.status_deploy_stopping(inputs)
});
/**
* | output |
* | --- |
* | "Updating" |
*
* @param {Status_Deploy_UpdatingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_deploy_updating = /** @type {((inputs?: Status_Deploy_UpdatingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Deploy_UpdatingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_deploy_updating(inputs)
	return __es.status_deploy_updating(inputs)
});
/**
* | output |
* | --- |
* | "Deployment" |
*
* @param {Status_DeploymentInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_deployment = /** @type {((inputs?: Status_DeploymentInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_DeploymentInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_deployment(inputs)
	return __es.status_deployment(inputs)
});
/**
* | output |
* | --- |
* | "Download Project" |
*
* @param {Status_Download_ProjectInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_download_project = /** @type {((inputs?: Status_Download_ProjectInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Download_ProjectInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_download_project(inputs)
	return __es.status_download_project(inputs)
});
/**
* | output |
* | --- |
* | "Error loading resources" |
*
* @param {Status_Error_Loading_ResourcesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_error_loading_resources = /** @type {((inputs?: Status_Error_Loading_ResourcesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Error_Loading_ResourcesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_error_loading_resources(inputs)
	return __es.status_error_loading_resources(inputs)
});
/**
* | output |
* | --- |
* | "Error loading tables" |
*
* @param {Status_Error_Loading_TablesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_error_loading_tables = /** @type {((inputs?: Status_Error_Loading_TablesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Error_Loading_TablesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_error_loading_tables(inputs)
	return __es.status_error_loading_tables(inputs)
});
/**
* | output |
* | --- |
* | "Errors" |
*
* @param {Status_ErrorsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_errors = /** @type {((inputs?: Status_ErrorsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_ErrorsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_errors(inputs)
	return __es.status_errors(inputs)
});
/**
* | output |
* | --- |
* | "External Tables" |
*
* @param {Status_External_Tables_SectionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_external_tables_section = /** @type {((inputs?: Status_External_Tables_SectionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_External_Tables_SectionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_external_tables_section(inputs)
	return __es.status_external_tables_section(inputs)
});
/**
* | output |
* | --- |
* | "Error" |
*
* @param {Status_Filter_ErrorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_filter_error = /** @type {((inputs?: Status_Filter_ErrorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Filter_ErrorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_filter_error(inputs)
	return __es.status_filter_error(inputs)
});
/**
* | output |
* | --- |
* | "OK" |
*
* @param {Status_Filter_OkInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_filter_ok = /** @type {((inputs?: Status_Filter_OkInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Filter_OkInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_filter_ok(inputs)
	return __es.status_filter_ok(inputs)
});
/**
* | output |
* | --- |
* | "Warn" |
*
* @param {Status_Filter_WarnInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_filter_warn = /** @type {((inputs?: Status_Filter_WarnInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Filter_WarnInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_filter_warn(inputs)
	return __es.status_filter_warn(inputs)
});
/**
* | output |
* | --- |
* | "⚠️ Warning: A full refresh will re-ingest ALL data from scratch. This operation can take a significant amount of time and will update all dependent resources..." |
*
* @param {Status_Full_Refresh_WarningInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_full_refresh_warning = /** @type {((inputs?: Status_Full_Refresh_WarningInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Full_Refresh_WarningInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_full_refresh_warning(inputs)
	return __es.status_full_refresh_warning(inputs)
});
/**
* | output |
* | --- |
* | "Refreshing this resource will update all dependent resources." |
*
* @param {Status_Incremental_Refresh_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_incremental_refresh_description = /** @type {((inputs?: Status_Incremental_Refresh_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Incremental_Refresh_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_incremental_refresh_description(inputs)
	return __es.status_incremental_refresh_description(inputs)
});
/**
* | output |
* | --- |
* | "AI Connector" |
*
* @param {Status_Label_Ai_ConnectorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_label_ai_connector = /** @type {((inputs?: Status_Label_Ai_ConnectorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Label_Ai_ConnectorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_label_ai_connector(inputs)
	return __es.status_label_ai_connector(inputs)
});
/**
* | output |
* | --- |
* | "Branch" |
*
* @param {Status_Label_BranchInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_label_branch = /** @type {((inputs?: Status_Label_BranchInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Label_BranchInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_label_branch(inputs)
	return __es.status_label_branch(inputs)
});
/**
* | output |
* | --- |
* | "Cluster Size" |
*
* @param {Status_Label_Cluster_SizeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_label_cluster_size = /** @type {((inputs?: Status_Label_Cluster_SizeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Label_Cluster_SizeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_label_cluster_size(inputs)
	return __es.status_label_cluster_size(inputs)
});
/**
* | output |
* | --- |
* | "Environment" |
*
* @param {Status_Label_EnvironmentInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_label_environment = /** @type {((inputs?: Status_Label_EnvironmentInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Label_EnvironmentInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_label_environment(inputs)
	return __es.status_label_environment(inputs)
});
/**
* | output |
* | --- |
* | "Last synced" |
*
* @param {Status_Label_Last_SyncedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_label_last_synced = /** @type {((inputs?: Status_Label_Last_SyncedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Label_Last_SyncedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_label_last_synced(inputs)
	return __es.status_label_last_synced(inputs)
});
/**
* | output |
* | --- |
* | "OLAP Engine" |
*
* @param {Status_Label_Olap_EngineInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_label_olap_engine = /** @type {((inputs?: Status_Label_Olap_EngineInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Label_Olap_EngineInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_label_olap_engine(inputs)
	return __es.status_label_olap_engine(inputs)
});
/**
* | output |
* | --- |
* | "Repo" |
*
* @param {Status_Label_RepoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_label_repo = /** @type {((inputs?: Status_Label_RepoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Label_RepoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_label_repo(inputs)
	return __es.status_label_repo(inputs)
});
/**
* | output |
* | --- |
* | "Runtime" |
*
* @param {Status_Label_RuntimeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_label_runtime = /** @type {((inputs?: Status_Label_RuntimeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Label_RuntimeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_label_runtime(inputs)
	return __es.status_label_runtime(inputs)
});
/**
* | output |
* | --- |
* | "Status" |
*
* @param {Status_Label_StatusInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_label_status = /** @type {((inputs?: Status_Label_StatusInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Label_StatusInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_label_status(inputs)
	return __es.status_label_status(inputs)
});
/**
* | output |
* | --- |
* | "Learn about connecting external OLAP engines" |
*
* @param {Status_Learn_About_External_OlapInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_learn_about_external_olap = /** @type {((inputs?: Status_Learn_About_External_OlapInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Learn_About_External_OlapInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_learn_about_external_olap(inputs)
	return __es.status_learn_about_external_olap(inputs)
});
/**
* | output |
* | --- |
* | "Learn more ->" |
*
* @param {Status_Learn_MoreInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_learn_more = /** @type {((inputs?: Status_Learn_MoreInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Learn_MoreInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_learn_more(inputs)
	return __es.status_learn_more(inputs)
});
/**
* | output |
* | --- |
* | "Load more tables" |
*
* @param {Status_Load_More_TablesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_load_more_tables = /** @type {((inputs?: Status_Load_More_TablesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Load_More_TablesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_load_more_tables(inputs)
	return __es.status_load_more_tables(inputs)
});
/**
* | output |
* | --- |
* | "Loading..." |
*
* @param {Status_LoadingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_loading = /** @type {((inputs?: Status_LoadingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_LoadingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_loading(inputs)
	return __es.status_loading(inputs)
});
/**
* | output |
* | --- |
* | "Loading models" |
*
* @param {Status_Loading_ModelsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_loading_models = /** @type {((inputs?: Status_Loading_ModelsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Loading_ModelsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_loading_models(inputs)
	return __es.status_loading_models(inputs)
});
/**
* | output |
* | --- |
* | "Loading resources..." |
*
* @param {Status_Loading_ResourcesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_loading_resources = /** @type {((inputs?: Status_Loading_ResourcesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Loading_ResourcesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_loading_resources(inputs)
	return __es.status_loading_resources(inputs)
});
/**
* | output |
* | --- |
* | "Loading tables..." |
*
* @param {Status_Loading_TablesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_loading_tables = /** @type {((inputs?: Status_Loading_TablesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Loading_TablesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_loading_tables(inputs)
	return __es.status_loading_tables(inputs)
});
/**
* | output |
* | --- |
* | "All levels" |
*
* @param {Status_Logs_All_LevelsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_all_levels = /** @type {((inputs?: Status_Logs_All_LevelsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_All_LevelsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_all_levels(inputs)
	return __es.status_logs_all_levels(inputs)
});
/**
* | output |
* | --- |
* | "Connecting" |
*
* @param {Status_Logs_ConnectingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_connecting = /** @type {((inputs?: Status_Logs_ConnectingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_ConnectingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_connecting(inputs)
	return __es.status_logs_connecting(inputs)
});
/**
* | output |
* | --- |
* | "Connection failed" |
*
* @param {Status_Logs_Connection_FailedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_connection_failed = /** @type {((inputs?: Status_Logs_Connection_FailedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_Connection_FailedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_connection_failed(inputs)
	return __es.status_logs_connection_failed(inputs)
});
/**
* | output |
* | --- |
* | "Disconnected" |
*
* @param {Status_Logs_DisconnectedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_disconnected = /** @type {((inputs?: Status_Logs_DisconnectedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_DisconnectedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_disconnected(inputs)
	return __es.status_logs_disconnected(inputs)
});
/**
* | output |
* | --- |
* | "Idle" |
*
* @param {Status_Logs_IdleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_idle = /** @type {((inputs?: Status_Logs_IdleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_IdleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_idle(inputs)
	return __es.status_logs_idle(inputs)
});
/**
* | output |
* | --- |
* | "Debug" |
*
* @param {Status_Logs_Level_DebugInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_level_debug = /** @type {((inputs?: Status_Logs_Level_DebugInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_Level_DebugInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_level_debug(inputs)
	return __es.status_logs_level_debug(inputs)
});
/**
* | output |
* | --- |
* | "Error" |
*
* @param {Status_Logs_Level_ErrorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_level_error = /** @type {((inputs?: Status_Logs_Level_ErrorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_Level_ErrorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_level_error(inputs)
	return __es.status_logs_level_error(inputs)
});
/**
* | output |
* | --- |
* | "Info" |
*
* @param {Status_Logs_Level_InfoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_level_info = /** @type {((inputs?: Status_Logs_Level_InfoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_Level_InfoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_level_info(inputs)
	return __es.status_logs_level_info(inputs)
});
/**
* | output |
* | --- |
* | "Warn" |
*
* @param {Status_Logs_Level_WarnInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_level_warn = /** @type {((inputs?: Status_Logs_Level_WarnInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_Level_WarnInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_level_warn(inputs)
	return __es.status_logs_level_warn(inputs)
});
/**
* | output |
* | --- |
* | "Live" |
*
* @param {Status_Logs_LiveInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_live = /** @type {((inputs?: Status_Logs_LiveInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_LiveInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_live(inputs)
	return __es.status_logs_live(inputs)
});
/**
* | output |
* | --- |
* | "No logs match the current filters" |
*
* @param {Status_Logs_No_MatchInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_no_match = /** @type {((inputs?: Status_Logs_No_MatchInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_No_MatchInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_no_match(inputs)
	return __es.status_logs_no_match(inputs)
});
/**
* | output |
* | --- |
* | "Retry" |
*
* @param {Status_Logs_RetryInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_retry = /** @type {((inputs?: Status_Logs_RetryInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_RetryInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_retry(inputs)
	return __es.status_logs_retry(inputs)
});
/**
* | output |
* | --- |
* | "Waiting for logs..." |
*
* @param {Status_Logs_WaitingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_logs_waiting = /** @type {((inputs?: Status_Logs_WaitingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Logs_WaitingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_logs_waiting(inputs)
	return __es.status_logs_waiting(inputs)
});
/**
* | output |
* | --- |
* | "Model Partitions" |
*
* @param {Status_Model_PartitionsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_model_partitions = /** @type {((inputs?: Status_Model_PartitionsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Model_PartitionsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_model_partitions(inputs)
	return __es.status_model_partitions(inputs)
});
/**
* | output |
* | --- |
* | "Models are created in Rill Developer." |
*
* @param {Status_Models_Created_In_DeveloperInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_models_created_in_developer = /** @type {((inputs?: Status_Models_Created_In_DeveloperInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Models_Created_In_DeveloperInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_models_created_in_developer(inputs)
	return __es.status_models_created_in_developer(inputs)
});
/**
* | output |
* | --- |
* | "Models" |
*
* @param {Status_Models_SectionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_models_section = /** @type {((inputs?: Status_Models_SectionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Models_SectionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_models_section(inputs)
	return __es.status_models_section(inputs)
});
/**
* | output |
* | --- |
* | "Analytics" |
*
* @param {Status_Nav_AnalyticsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_nav_analytics = /** @type {((inputs?: Status_Nav_AnalyticsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Nav_AnalyticsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_nav_analytics(inputs)
	return __es.status_nav_analytics(inputs)
});
/**
* | output |
* | --- |
* | "Branches" |
*
* @param {Status_Nav_BranchesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_nav_branches = /** @type {((inputs?: Status_Nav_BranchesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Nav_BranchesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_nav_branches(inputs)
	return __es.status_nav_branches(inputs)
});
/**
* | output |
* | --- |
* | "Logs" |
*
* @param {Status_Nav_LogsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_nav_logs = /** @type {((inputs?: Status_Nav_LogsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Nav_LogsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_nav_logs(inputs)
	return __es.status_nav_logs(inputs)
});
/**
* | output |
* | --- |
* | "Overview" |
*
* @param {Status_Nav_OverviewInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_nav_overview = /** @type {((inputs?: Status_Nav_OverviewInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Nav_OverviewInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_nav_overview(inputs)
	return __es.status_nav_overview(inputs)
});
/**
* | output |
* | --- |
* | "Resources" |
*
* @param {Status_Nav_ResourcesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_nav_resources = /** @type {((inputs?: Status_Nav_ResourcesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Nav_ResourcesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_nav_resources(inputs)
	return __es.status_nav_resources(inputs)
});
/**
* | output |
* | --- |
* | "Tables" |
*
* @param {Status_Nav_TablesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_nav_tables = /** @type {((inputs?: Status_Nav_TablesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Nav_TablesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_nav_tables(inputs)
	return __es.status_nav_tables(inputs)
});
/**
* | output |
* | --- |
* | "No errors detected." |
*
* @param {Status_No_ErrorsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_no_errors = /** @type {((inputs?: Status_No_ErrorsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_No_ErrorsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_no_errors(inputs)
	return __es.status_no_errors(inputs)
});
/**
* | output |
* | --- |
* | "No external tables" |
*
* @param {Status_No_External_TablesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_no_external_tables = /** @type {((inputs?: Status_No_External_TablesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_No_External_TablesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_no_external_tables(inputs)
	return __es.status_no_external_tables(inputs)
});
/**
* | output |
* | --- |
* | "No external tables match the current filters" |
*
* @param {Status_No_External_Tables_Match_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_no_external_tables_match_filters = /** @type {((inputs?: Status_No_External_Tables_Match_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_No_External_Tables_Match_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_no_external_tables_match_filters(inputs)
	return __es.status_no_external_tables_match_filters(inputs)
});
/**
* | output |
* | --- |
* | "No models" |
*
* @param {Status_No_ModelsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_no_models = /** @type {((inputs?: Status_No_ModelsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_No_ModelsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_no_models(inputs)
	return __es.status_no_models(inputs)
});
/**
* | output |
* | --- |
* | "No models match the current filters" |
*
* @param {Status_No_Models_Match_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_no_models_match_filters = /** @type {((inputs?: Status_No_Models_Match_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_No_Models_Match_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_no_models_match_filters(inputs)
	return __es.status_no_models_match_filters(inputs)
});
/**
* | output |
* | --- |
* | "No parse errors" |
*
* @param {Status_No_Parse_ErrorsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_no_parse_errors = /** @type {((inputs?: Status_No_Parse_ErrorsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_No_Parse_ErrorsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_no_parse_errors(inputs)
	return __es.status_no_parse_errors(inputs)
});
/**
* | output |
* | --- |
* | "No resource data available" |
*
* @param {Status_No_Resource_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_no_resource_data = /** @type {((inputs?: Status_No_Resource_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_No_Resource_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_no_resource_data(inputs)
	return __es.status_no_resource_data(inputs)
});
/**
* | output |
* | --- |
* | "No resources found." |
*
* @param {Status_No_ResourcesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_no_resources = /** @type {((inputs?: Status_No_ResourcesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_No_ResourcesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_no_resources(inputs)
	return __es.status_no_resources(inputs)
});
/**
* | output |
* | --- |
* | "No resources match the current filters" |
*
* @param {Status_No_Resources_Match_FiltersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_no_resources_match_filters = /** @type {((inputs?: Status_No_Resources_Match_FiltersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_No_Resources_Match_FiltersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_no_resources_match_filters(inputs)
	return __es.status_no_resources_match_filters(inputs)
});
/**
* | output |
* | --- |
* | "No tables found." |
*
* @param {Status_No_TablesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_no_tables = /** @type {((inputs?: Status_No_TablesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_No_TablesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_no_tables(inputs)
	return __es.status_no_tables(inputs)
});
/**
* | output |
* | --- |
* | "Note:" |
*
* @param {Status_NoteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_note = /** @type {((inputs?: Status_NoteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_NoteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_note(inputs)
	return __es.status_note(inputs)
});
/**
* | output |
* | --- |
* | "Owned by another user" |
*
* @param {Status_Owned_By_OtherInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_owned_by_other = /** @type {((inputs?: Status_Owned_By_OtherInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Owned_By_OtherInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_owned_by_other(inputs)
	return __es.status_owned_by_other(inputs)
});
/**
* | output |
* | --- |
* | "Owned by you" |
*
* @param {Status_Owned_By_YouInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_owned_by_you = /** @type {((inputs?: Status_Owned_By_YouInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Owned_By_YouInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_owned_by_you(inputs)
	return __es.status_owned_by_you(inputs)
});
/**
* | output |
* | --- |
* | "Project Status" |
*
* @param {Status_Page_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_page_title = /** @type {((inputs?: Status_Page_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Page_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_page_title(inputs)
	return __es.status_page_title(inputs)
});
/**
* | output |
* | --- |
* | "Parse error" |
*
* @param {Status_Parse_ErrorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_parse_error = /** @type {((inputs?: Status_Parse_ErrorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Parse_ErrorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_parse_error(inputs)
	return __es.status_parse_error(inputs)
});
/**
* | output |
* | --- |
* | "Parse errors" |
*
* @param {Status_Parse_ErrorsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_parse_errors = /** @type {((inputs?: Status_Parse_ErrorsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Parse_ErrorsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_parse_errors(inputs)
	return __es.status_parse_errors(inputs)
});
/**
* | output |
* | --- |
* | "Parse Errors" |
*
* @param {Status_Parse_Errors_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_parse_errors_title = /** @type {((inputs?: Status_Parse_Errors_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Parse_Errors_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_parse_errors_title(inputs)
	return __es.status_parse_errors_title(inputs)
});
/**
* | output |
* | --- |
* | "Reconciling" |
*
* @param {Status_ReconcilingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_reconciling = /** @type {((inputs?: Status_ReconcilingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_ReconcilingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_reconciling(inputs)
	return __es.status_reconciling(inputs)
});
/**
* | output |
* | --- |
* | "Refresh all" |
*
* @param {Status_Refresh_AllInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_refresh_all = /** @type {((inputs?: Status_Refresh_AllInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Refresh_AllInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_refresh_all(inputs)
	return __es.status_refresh_all(inputs)
});
/**
* | output |
* | --- |
* | "This will refresh all project sources and models." |
*
* @param {Status_Refresh_All_Confirm_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_refresh_all_confirm_body = /** @type {((inputs?: Status_Refresh_All_Confirm_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Refresh_All_Confirm_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_refresh_all_confirm_body(inputs)
	return __es.status_refresh_all_confirm_body(inputs)
});
/**
* | output |
* | --- |
* | "To refresh a single resource, scroll to the source or model, click the '...' button, and select the refresh option." |
*
* @param {Status_Refresh_All_Confirm_TipInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_refresh_all_confirm_tip = /** @type {((inputs?: Status_Refresh_All_Confirm_TipInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Refresh_All_Confirm_TipInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_refresh_all_confirm_tip(inputs)
	return __es.status_refresh_all_confirm_tip(inputs)
});
/**
* | output |
* | --- |
* | "Refresh all sources and models?" |
*
* @param {Status_Refresh_All_Confirm_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_refresh_all_confirm_title = /** @type {((inputs?: Status_Refresh_All_Confirm_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Refresh_All_Confirm_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_refresh_all_confirm_title(inputs)
	return __es.status_refresh_all_confirm_title(inputs)
});
/**
* | output |
* | --- |
* | "Refresh all sources and models" |
*
* @param {Status_Refresh_All_Sources_ModelsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_refresh_all_sources_models = /** @type {((inputs?: Status_Refresh_All_Sources_ModelsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Refresh_All_Sources_ModelsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_refresh_all_sources_models(inputs)
	return __es.status_refresh_all_sources_models(inputs)
});
/**
* | output |
* | --- |
* | "This will re-execute all partitions that failed during their last run. The refresh will happen in the background." |
*
* @param {Status_Refresh_Errored_Confirm_BodyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_refresh_errored_confirm_body = /** @type {((inputs?: Status_Refresh_Errored_Confirm_BodyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Refresh_Errored_Confirm_BodyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_refresh_errored_confirm_body(inputs)
	return __es.status_refresh_errored_confirm_body(inputs)
});
/**
* | output |
* | --- |
* | "Refresh Errored Partitions for {modelName}?" |
*
* @param {Status_Refresh_Errored_Confirm_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_refresh_errored_confirm_title = /** @type {((inputs: Status_Refresh_Errored_Confirm_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Refresh_Errored_Confirm_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_refresh_errored_confirm_title(inputs)
	return __es.status_refresh_errored_confirm_title(inputs)
});
/**
* | output |
* | --- |
* | "Refreshing..." |
*
* @param {Status_RefreshingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_refreshing = /** @type {((inputs?: Status_RefreshingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_RefreshingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_refreshing(inputs)
	return __es.status_refreshing(inputs)
});
/**
* | output |
* | --- |
* | "Resource is currently being reconciled" |
*
* @param {Status_Resource_ReconcilingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_resource_reconciling = /** @type {((inputs?: Status_Resource_ReconcilingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Resource_ReconcilingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_resource_reconciling(inputs)
	return __es.status_resource_reconciling(inputs)
});
/**
* | output |
* | --- |
* | "Rill Managed" |
*
* @param {Status_Rill_ManagedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_rill_managed = /** @type {((inputs?: Status_Rill_ManagedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Rill_ManagedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_rill_managed(inputs)
	return __es.status_rill_managed(inputs)
});
/**
* | output |
* | --- |
* | "Tables" |
*
* @param {Status_Table_PluralInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_table_plural = /** @type {((inputs?: Status_Table_PluralInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Table_PluralInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_table_plural(inputs)
	return __es.status_table_plural(inputs)
});
/**
* | output |
* | --- |
* | "Table" |
*
* @param {Status_Table_SingularInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_table_singular = /** @type {((inputs?: Status_Table_SingularInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Table_SingularInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_table_singular(inputs)
	return __es.status_table_singular(inputs)
});
/**
* | output |
* | --- |
* | "Unable to check for errors." |
*
* @param {Status_Unable_To_Check_ErrorsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_unable_to_check_errors = /** @type {((inputs?: Status_Unable_To_Check_ErrorsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Unable_To_Check_ErrorsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_unable_to_check_errors(inputs)
	return __es.status_unable_to_check_errors(inputs)
});
/**
* | output |
* | --- |
* | "View all" |
*
* @param {Status_View_AllInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_view_all = /** @type {((inputs?: Status_View_AllInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_View_AllInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_view_all(inputs)
	return __es.status_view_all(inputs)
});
/**
* | output |
* | --- |
* | "Views" |
*
* @param {Status_View_PluralInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_view_plural = /** @type {((inputs?: Status_View_PluralInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_View_PluralInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_view_plural(inputs)
	return __es.status_view_plural(inputs)
});
/**
* | output |
* | --- |
* | "View" |
*
* @param {Status_View_SingularInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_view_singular = /** @type {((inputs?: Status_View_SingularInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_View_SingularInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_view_singular(inputs)
	return __es.status_view_singular(inputs)
});
/**
* | output |
* | --- |
* | "Yes, refresh" |
*
* @param {Status_Yes_RefreshInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const status_yes_refresh = /** @type {((inputs?: Status_Yes_RefreshInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Status_Yes_RefreshInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.status_yes_refresh(inputs)
	return __es.status_yes_refresh(inputs)
});
/**
* | output |
* | --- |
* | "Dark" |
*
* @param {Theme_DarkInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const theme_dark = /** @type {((inputs?: Theme_DarkInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Theme_DarkInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.theme_dark(inputs)
	return __es.theme_dark(inputs)
});
/**
* | output |
* | --- |
* | "Theme" |
*
* @param {Theme_LabelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const theme_label = /** @type {((inputs?: Theme_LabelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Theme_LabelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.theme_label(inputs)
	return __es.theme_label(inputs)
});
/**
* | output |
* | --- |
* | "Light" |
*
* @param {Theme_LightInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const theme_light = /** @type {((inputs?: Theme_LightInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Theme_LightInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.theme_light(inputs)
	return __es.theme_light(inputs)
});
/**
* | output |
* | --- |
* | "System" |
*
* @param {Theme_SystemInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const theme_system = /** @type {((inputs?: Theme_SystemInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Theme_SystemInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.theme_system(inputs)
	return __es.theme_system(inputs)
});
/**
* | output |
* | --- |
* | "1 day ago" |
*
* @param {Time_1_Day_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_1_day_ago = /** @type {((inputs?: Time_1_Day_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_1_Day_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_1_day_ago(inputs)
	return __es.time_1_day_ago(inputs)
});
/**
* | output |
* | --- |
* | "1 hour ago" |
*
* @param {Time_1_Hour_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_1_hour_ago = /** @type {((inputs?: Time_1_Hour_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_1_Hour_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_1_hour_ago(inputs)
	return __es.time_1_hour_ago(inputs)
});
/**
* | output |
* | --- |
* | "1 minute ago" |
*
* @param {Time_1_Minute_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_1_minute_ago = /** @type {((inputs?: Time_1_Minute_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_1_Minute_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_1_minute_ago(inputs)
	return __es.time_1_minute_ago(inputs)
});
/**
* | output |
* | --- |
* | "1 month ago" |
*
* @param {Time_1_Month_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_1_month_ago = /** @type {((inputs?: Time_1_Month_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_1_Month_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_1_month_ago(inputs)
	return __es.time_1_month_ago(inputs)
});
/**
* | output |
* | --- |
* | "1 week ago" |
*
* @param {Time_1_Week_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_1_week_ago = /** @type {((inputs?: Time_1_Week_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_1_Week_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_1_week_ago(inputs)
	return __es.time_1_week_ago(inputs)
});
/**
* | output |
* | --- |
* | "1 year ago" |
*
* @param {Time_1_Year_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_1_year_ago = /** @type {((inputs?: Time_1_Year_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_1_Year_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_1_year_ago(inputs)
	return __es.time_1_year_ago(inputs)
});
/**
* | output |
* | --- |
* | "{duration} ago" |
*
* @param {Time_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_ago = /** @type {((inputs: Time_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_ago(inputs)
	return __es.time_ago(inputs)
});
/**
* | output |
* | --- |
* | "All Time" |
*
* @param {Time_All_TimeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_all_time = /** @type {((inputs?: Time_All_TimeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_All_TimeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_all_time(inputs)
	return __es.time_all_time(inputs)
});
/**
* | output |
* | --- |
* | "Comparing" |
*
* @param {Time_ComparingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_comparing = /** @type {((inputs?: Time_ComparingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_ComparingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_comparing(inputs)
	return __es.time_comparing(inputs)
});
/**
* | output |
* | --- |
* | "Previous day" |
*
* @param {Time_Comparison_Previous_DayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_comparison_previous_day = /** @type {((inputs?: Time_Comparison_Previous_DayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Comparison_Previous_DayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_comparison_previous_day(inputs)
	return __es.time_comparison_previous_day(inputs)
});
/**
* | output |
* | --- |
* | "Previous period" |
*
* @param {Time_Comparison_Previous_PeriodInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_comparison_previous_period = /** @type {((inputs?: Time_Comparison_Previous_PeriodInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Comparison_Previous_PeriodInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_comparison_previous_period(inputs)
	return __es.time_comparison_previous_period(inputs)
});
/**
* | output |
* | --- |
* | "Custom" |
*
* @param {Time_CustomInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_custom = /** @type {((inputs?: Time_CustomInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_CustomInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_custom(inputs)
	return __es.time_custom(inputs)
});
/**
* | output |
* | --- |
* | "Custom range" |
*
* @param {Time_Custom_RangeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_custom_range = /** @type {((inputs?: Time_Custom_RangeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Custom_RangeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_custom_range(inputs)
	return __es.time_custom_range(inputs)
});
/**
* | output |
* | --- |
* | "Enter a time range" |
*
* @param {Time_Enter_Time_RangeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_enter_time_range = /** @type {((inputs?: Time_Enter_Time_RangeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Enter_Time_RangeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_enter_time_range(inputs)
	return __es.time_enter_time_range(inputs)
});
/**
* | output |
* | --- |
* | "{duration} from now" |
*
* @param {Time_From_NowInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_from_now = /** @type {((inputs: Time_From_NowInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_From_NowInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_from_now(inputs)
	return __es.time_from_now(inputs)
});
/**
* | output |
* | --- |
* | "by" |
*
* @param {Time_Grain_ByInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_by = /** @type {((inputs?: Time_Grain_ByInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_ByInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_by(inputs)
	return __es.time_grain_by(inputs)
});
/**
* | output |
* | --- |
* | "complete" |
*
* @param {Time_Grain_CompleteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_complete = /** @type {((inputs?: Time_Grain_CompleteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_CompleteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_complete(inputs)
	return __es.time_grain_complete(inputs)
});
/**
* | output |
* | --- |
* | "day" |
*
* @param {Time_Grain_DayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_day = /** @type {((inputs?: Time_Grain_DayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_DayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_day(inputs)
	return __es.time_grain_day(inputs)
});
/**
* | output |
* | --- |
* | "days" |
*
* @param {Time_Grain_DaysInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_days = /** @type {((inputs?: Time_Grain_DaysInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_DaysInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_days(inputs)
	return __es.time_grain_days(inputs)
});
/**
* | output |
* | --- |
* | "hour" |
*
* @param {Time_Grain_HourInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_hour = /** @type {((inputs?: Time_Grain_HourInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_HourInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_hour(inputs)
	return __es.time_grain_hour(inputs)
});
/**
* | output |
* | --- |
* | "hours" |
*
* @param {Time_Grain_HoursInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_hours = /** @type {((inputs?: Time_Grain_HoursInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_HoursInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_hours(inputs)
	return __es.time_grain_hours(inputs)
});
/**
* | output |
* | --- |
* | "minute" |
*
* @param {Time_Grain_MinuteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_minute = /** @type {((inputs?: Time_Grain_MinuteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_MinuteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_minute(inputs)
	return __es.time_grain_minute(inputs)
});
/**
* | output |
* | --- |
* | "minutes" |
*
* @param {Time_Grain_MinutesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_minutes = /** @type {((inputs?: Time_Grain_MinutesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_MinutesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_minutes(inputs)
	return __es.time_grain_minutes(inputs)
});
/**
* | output |
* | --- |
* | "month" |
*
* @param {Time_Grain_MonthInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_month = /** @type {((inputs?: Time_Grain_MonthInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_MonthInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_month(inputs)
	return __es.time_grain_month(inputs)
});
/**
* | output |
* | --- |
* | "months" |
*
* @param {Time_Grain_MonthsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_months = /** @type {((inputs?: Time_Grain_MonthsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_MonthsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_months(inputs)
	return __es.time_grain_months(inputs)
});
/**
* | output |
* | --- |
* | "quarter" |
*
* @param {Time_Grain_QuarterInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_quarter = /** @type {((inputs?: Time_Grain_QuarterInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_QuarterInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_quarter(inputs)
	return __es.time_grain_quarter(inputs)
});
/**
* | output |
* | --- |
* | "quarters" |
*
* @param {Time_Grain_QuartersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_quarters = /** @type {((inputs?: Time_Grain_QuartersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_QuartersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_quarters(inputs)
	return __es.time_grain_quarters(inputs)
});
/**
* | output |
* | --- |
* | "Time" |
*
* @param {Time_Grain_TimeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_time = /** @type {((inputs?: Time_Grain_TimeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_TimeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_time(inputs)
	return __es.time_grain_time(inputs)
});
/**
* | output |
* | --- |
* | "week" |
*
* @param {Time_Grain_WeekInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_week = /** @type {((inputs?: Time_Grain_WeekInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_WeekInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_week(inputs)
	return __es.time_grain_week(inputs)
});
/**
* | output |
* | --- |
* | "weeks" |
*
* @param {Time_Grain_WeeksInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_weeks = /** @type {((inputs?: Time_Grain_WeeksInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_WeeksInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_weeks(inputs)
	return __es.time_grain_weeks(inputs)
});
/**
* | output |
* | --- |
* | "year" |
*
* @param {Time_Grain_YearInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_year = /** @type {((inputs?: Time_Grain_YearInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_YearInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_year(inputs)
	return __es.time_grain_year(inputs)
});
/**
* | output |
* | --- |
* | "years" |
*
* @param {Time_Grain_YearsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_grain_years = /** @type {((inputs?: Time_Grain_YearsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Grain_YearsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_grain_years(inputs)
	return __es.time_grain_years(inputs)
});
/**
* | output |
* | --- |
* | "Just now" |
*
* @param {Time_Just_NowInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_just_now = /** @type {((inputs?: Time_Just_NowInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Just_NowInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_just_now(inputs)
	return __es.time_just_now(inputs)
});
/**
* | output |
* | --- |
* | "Last {duration}" |
*
* @param {Time_Last_DurationInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_last_duration = /** @type {((inputs: Time_Last_DurationInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Last_DurationInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_last_duration(inputs)
	return __es.time_last_duration(inputs)
});
/**
* | output |
* | --- |
* | "Month to Date" |
*
* @param {Time_Month_To_DateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_month_to_date = /** @type {((inputs?: Time_Month_To_DateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Month_To_DateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_month_to_date(inputs)
	return __es.time_month_to_date(inputs)
});
/**
* | output |
* | --- |
* | "{count} days ago" |
*
* @param {Time_N_Days_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_n_days_ago = /** @type {((inputs: Time_N_Days_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_N_Days_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_n_days_ago(inputs)
	return __es.time_n_days_ago(inputs)
});
/**
* | output |
* | --- |
* | "{count} hours ago" |
*
* @param {Time_N_Hours_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_n_hours_ago = /** @type {((inputs: Time_N_Hours_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_N_Hours_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_n_hours_ago(inputs)
	return __es.time_n_hours_ago(inputs)
});
/**
* | output |
* | --- |
* | "{count} minutes ago" |
*
* @param {Time_N_Minutes_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_n_minutes_ago = /** @type {((inputs: Time_N_Minutes_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_N_Minutes_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_n_minutes_ago(inputs)
	return __es.time_n_minutes_ago(inputs)
});
/**
* | output |
* | --- |
* | "{count} months ago" |
*
* @param {Time_N_Months_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_n_months_ago = /** @type {((inputs: Time_N_Months_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_N_Months_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_n_months_ago(inputs)
	return __es.time_n_months_ago(inputs)
});
/**
* | output |
* | --- |
* | "{count} weeks ago" |
*
* @param {Time_N_Weeks_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_n_weeks_ago = /** @type {((inputs: Time_N_Weeks_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_N_Weeks_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_n_weeks_ago(inputs)
	return __es.time_n_weeks_ago(inputs)
});
/**
* | output |
* | --- |
* | "{count} years ago" |
*
* @param {Time_N_Years_AgoInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_n_years_ago = /** @type {((inputs: Time_N_Years_AgoInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_N_Years_AgoInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_n_years_ago(inputs)
	return __es.time_n_years_ago(inputs)
});
/**
* | output |
* | --- |
* | "No comparison dimension" |
*
* @param {Time_No_Comparison_DimensionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_no_comparison_dimension = /** @type {((inputs?: Time_No_Comparison_DimensionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_No_Comparison_DimensionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_no_comparison_dimension(inputs)
	return __es.time_no_comparison_dimension(inputs)
});
/**
* | output |
* | --- |
* | "no comparison period" |
*
* @param {Time_No_Comparison_PeriodInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_no_comparison_period = /** @type {((inputs?: Time_No_Comparison_PeriodInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_No_Comparison_PeriodInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_no_comparison_period(inputs)
	return __es.time_no_comparison_period(inputs)
});
/**
* | output |
* | --- |
* | "Previous month" |
*
* @param {Time_Previous_MonthInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_previous_month = /** @type {((inputs?: Time_Previous_MonthInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Previous_MonthInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_previous_month(inputs)
	return __es.time_previous_month(inputs)
});
/**
* | output |
* | --- |
* | "Previous quarter" |
*
* @param {Time_Previous_QuarterInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_previous_quarter = /** @type {((inputs?: Time_Previous_QuarterInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Previous_QuarterInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_previous_quarter(inputs)
	return __es.time_previous_quarter(inputs)
});
/**
* | output |
* | --- |
* | "Previous week" |
*
* @param {Time_Previous_WeekInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_previous_week = /** @type {((inputs?: Time_Previous_WeekInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Previous_WeekInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_previous_week(inputs)
	return __es.time_previous_week(inputs)
});
/**
* | output |
* | --- |
* | "Previous year" |
*
* @param {Time_Previous_YearInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_previous_year = /** @type {((inputs?: Time_Previous_YearInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Previous_YearInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_previous_year(inputs)
	return __es.time_previous_year(inputs)
});
/**
* | output |
* | --- |
* | "Quarter to Date" |
*
* @param {Time_Quarter_To_DateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_quarter_to_date = /** @type {((inputs?: Time_Quarter_To_DateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Quarter_To_DateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_quarter_to_date(inputs)
	return __es.time_quarter_to_date(inputs)
});
/**
* | output |
* | --- |
* | "{grain} to date" |
*
* @param {Time_Range_Grain_To_DateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_range_grain_to_date = /** @type {((inputs: Time_Range_Grain_To_DateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Range_Grain_To_DateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_range_grain_to_date(inputs)
	return __es.time_range_grain_to_date(inputs)
});
/**
* | output |
* | --- |
* | "last {count} {grains}" |
*
* @param {Time_Range_Last_N_GrainsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_range_last_n_grains = /** @type {((inputs: Time_Range_Last_N_GrainsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Range_Last_N_GrainsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_range_last_n_grains(inputs)
	return __es.time_range_last_n_grains(inputs)
});
/**
* | output |
* | --- |
* | "next {grain}" |
*
* @param {Time_Range_Next_GrainInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_range_next_grain = /** @type {((inputs: Time_Range_Next_GrainInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Range_Next_GrainInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_range_next_grain(inputs)
	return __es.time_range_next_grain(inputs)
});
/**
* | output |
* | --- |
* | "next {count} {grains}" |
*
* @param {Time_Range_Next_N_GrainsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_range_next_n_grains = /** @type {((inputs: Time_Range_Next_N_GrainsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Range_Next_N_GrainsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_range_next_n_grains(inputs)
	return __es.time_range_next_n_grains(inputs)
});
/**
* | output |
* | --- |
* | "previous {grain}" |
*
* @param {Time_Range_Previous_GrainInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_range_previous_grain = /** @type {((inputs: Time_Range_Previous_GrainInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Range_Previous_GrainInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_range_previous_grain(inputs)
	return __es.time_range_previous_grain(inputs)
});
/**
* | output |
* | --- |
* | "this {grain}" |
*
* @param {Time_Range_This_GrainInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_range_this_grain = /** @type {((inputs: Time_Range_This_GrainInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Range_This_GrainInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_range_this_grain(inputs)
	return __es.time_range_this_grain(inputs)
});
/**
* | output |
* | --- |
* | "complete" |
*
* @param {Time_Ref_CompleteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_ref_complete = /** @type {((inputs?: Time_Ref_CompleteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Ref_CompleteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_ref_complete(inputs)
	return __es.time_ref_complete(inputs)
});
/**
* | output |
* | --- |
* | "complete data" |
*
* @param {Time_Ref_Complete_DataInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_ref_complete_data = /** @type {((inputs?: Time_Ref_Complete_DataInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Ref_Complete_DataInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_ref_complete_data(inputs)
	return __es.time_ref_complete_data(inputs)
});
/**
* | output |
* | --- |
* | "current" |
*
* @param {Time_Ref_CurrentInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_ref_current = /** @type {((inputs?: Time_Ref_CurrentInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Ref_CurrentInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_ref_current(inputs)
	return __es.time_ref_current(inputs)
});
/**
* | output |
* | --- |
* | "earliest" |
*
* @param {Time_Ref_EarliestInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_ref_earliest = /** @type {((inputs?: Time_Ref_EarliestInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Ref_EarliestInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_ref_earliest(inputs)
	return __es.time_ref_earliest(inputs)
});
/**
* | output |
* | --- |
* | "latest" |
*
* @param {Time_Ref_LatestInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_ref_latest = /** @type {((inputs?: Time_Ref_LatestInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Ref_LatestInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_ref_latest(inputs)
	return __es.time_ref_latest(inputs)
});
/**
* | output |
* | --- |
* | "now" |
*
* @param {Time_Ref_NowInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_ref_now = /** @type {((inputs?: Time_Ref_NowInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Ref_NowInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_ref_now(inputs)
	return __es.time_ref_now(inputs)
});
/**
* | output |
* | --- |
* | "{count}h ago" |
*
* @param {Time_Relative_Hours_ShortInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_relative_hours_short = /** @type {((inputs: Time_Relative_Hours_ShortInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Relative_Hours_ShortInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_relative_hours_short(inputs)
	return __es.time_relative_hours_short(inputs)
});
/**
* | output |
* | --- |
* | "{count}m ago" |
*
* @param {Time_Relative_Minutes_ShortInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_relative_minutes_short = /** @type {((inputs: Time_Relative_Minutes_ShortInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Relative_Minutes_ShortInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_relative_minutes_short(inputs)
	return __es.time_relative_minutes_short(inputs)
});
/**
* | output |
* | --- |
* | "now" |
*
* @param {Time_Relative_NowInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_relative_now = /** @type {((inputs?: Time_Relative_NowInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Relative_NowInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_relative_now(inputs)
	return __es.time_relative_now(inputs)
});
/**
* | output |
* | --- |
* | "Today" |
*
* @param {Time_TodayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_today = /** @type {((inputs?: Time_TodayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_TodayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_today(inputs)
	return __es.time_today(inputs)
});
/**
* | output |
* | --- |
* | "Unable to parse time string" |
*
* @param {Time_Unable_To_ParseInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_unable_to_parse = /** @type {((inputs?: Time_Unable_To_ParseInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Unable_To_ParseInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_unable_to_parse(inputs)
	return __es.time_unable_to_parse(inputs)
});
/**
* | output |
* | --- |
* | "vs" |
*
* @param {Time_VsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_vs = /** @type {((inputs?: Time_VsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_VsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_vs(inputs)
	return __es.time_vs(inputs)
});
/**
* | output |
* | --- |
* | "Week to Date" |
*
* @param {Time_Week_To_DateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_week_to_date = /** @type {((inputs?: Time_Week_To_DateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Week_To_DateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_week_to_date(inputs)
	return __es.time_week_to_date(inputs)
});
/**
* | output |
* | --- |
* | "Year to Date" |
*
* @param {Time_Year_To_DateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_year_to_date = /** @type {((inputs?: Time_Year_To_DateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_Year_To_DateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_year_to_date(inputs)
	return __es.time_year_to_date(inputs)
});
/**
* | output |
* | --- |
* | "Yesterday" |
*
* @param {Time_YesterdayInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const time_yesterday = /** @type {((inputs?: Time_YesterdayInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Time_YesterdayInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.time_yesterday(inputs)
	return __es.time_yesterday(inputs)
});
/**
* | output |
* | --- |
* | "Project access changed to everyone" |
*
* @param {Users_Access_Changed_EveryoneInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_access_changed_everyone = /** @type {((inputs?: Users_Access_Changed_EveryoneInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Access_Changed_EveryoneInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_access_changed_everyone(inputs)
	return __es.users_access_changed_everyone(inputs)
});
/**
* | output |
* | --- |
* | "Project access changed to invite-only" |
*
* @param {Users_Access_Changed_Invite_OnlyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_access_changed_invite_only = /** @type {((inputs?: Users_Access_Changed_Invite_OnlyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Access_Changed_Invite_OnlyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_access_changed_invite_only(inputs)
	return __es.users_access_changed_invite_only(inputs)
});
/**
* | output |
* | --- |
* | "Access level" |
*
* @param {Users_Access_LevelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_access_level = /** @type {((inputs?: Users_Access_LevelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Access_LevelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_access_level(inputs)
	return __es.users_access_level(inputs)
});
/**
* | output |
* | --- |
* | "Add guest" |
*
* @param {Users_Add_GuestInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_add_guest = /** @type {((inputs?: Users_Add_GuestInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Add_GuestInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_add_guest(inputs)
	return __es.users_add_guest(inputs)
});
/**
* | output |
* | --- |
* | "Add guests" |
*
* @param {Users_Add_Guests_ButtonInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_add_guests_button = /** @type {((inputs?: Users_Add_Guests_ButtonInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Add_Guests_ButtonInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_add_guests_button(inputs)
	return __es.users_add_guests_button(inputs)
});
/**
* | output |
* | --- |
* | "Guests can only access provisioned projects with assigned roles. They do not have organization-wide access." |
*
* @param {Users_Add_Guests_DescriptionInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_add_guests_description = /** @type {((inputs?: Users_Add_Guests_DescriptionInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Add_Guests_DescriptionInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_add_guests_description(inputs)
	return __es.users_add_guests_description(inputs)
});
/**
* | output |
* | --- |
* | "Add guest users" |
*
* @param {Users_Add_Guests_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_add_guests_title = /** @type {((inputs?: Users_Add_Guests_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Add_Guests_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_add_guests_title(inputs)
	return __es.users_add_guests_title(inputs)
});
/**
* | output |
* | --- |
* | "Add users" |
*
* @param {Users_Add_UsersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_add_users = /** @type {((inputs?: Users_Add_UsersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Add_UsersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_add_users(inputs)
	return __es.users_add_users(inputs)
});
/**
* | output |
* | --- |
* | "Added {count} groups" |
*
* @param {Users_Added_Groups_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_added_groups_count = /** @type {((inputs: Users_Added_Groups_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Added_Groups_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_added_groups_count(inputs)
	return __es.users_added_groups_count(inputs)
});
/**
* | output |
* | --- |
* | "{emails} already a member of this organization" |
*
* @param {Users_Already_MemberInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_already_member = /** @type {((inputs: Users_Already_MemberInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Already_MemberInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_already_member(inputs)
	return __es.users_already_member(inputs)
});
/**
* | output |
* | --- |
* | "and {count} more" |
*
* @param {Users_And_MoreInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_and_more = /** @type {((inputs: Users_And_MoreInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_And_MoreInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_and_more(inputs)
	return __es.users_and_more(inputs)
});
/**
* | output |
* | --- |
* | "This user is currently the billing contact. To change their role, assign another admin as the billing contact first." |
*
* @param {Users_Billing_Change_Role_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_billing_change_role_desc = /** @type {((inputs?: Users_Billing_Change_Role_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Billing_Change_Role_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_billing_change_role_desc(inputs)
	return __es.users_billing_change_role_desc(inputs)
});
/**
* | output |
* | --- |
* | "Assign a new billing contact first to change this user's role" |
*
* @param {Users_Billing_Change_Role_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_billing_change_role_title = /** @type {((inputs?: Users_Billing_Change_Role_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Billing_Change_Role_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_billing_change_role_title(inputs)
	return __es.users_billing_change_role_title(inputs)
});
/**
* | output |
* | --- |
* | "This user is the current billing contact and can't be removed. To proceed, assign another admin as the billing contact." |
*
* @param {Users_Billing_Remove_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_billing_remove_desc = /** @type {((inputs?: Users_Billing_Remove_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Billing_Remove_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_billing_remove_desc(inputs)
	return __es.users_billing_remove_desc(inputs)
});
/**
* | output |
* | --- |
* | "Assign a new billing contact first to remove this user" |
*
* @param {Users_Billing_Remove_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_billing_remove_title = /** @type {((inputs?: Users_Billing_Remove_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Billing_Remove_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_billing_remove_title(inputs)
	return __es.users_billing_remove_title(inputs)
});
/**
* | output |
* | --- |
* | "Cancel" |
*
* @param {Users_CancelInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_cancel = /** @type {((inputs?: Users_CancelInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_CancelInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_cancel(inputs)
	return __es.users_cancel(inputs)
});
/**
* | output |
* | --- |
* | "Change billing contact" |
*
* @param {Users_Change_Billing_ContactInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_change_billing_contact = /** @type {((inputs?: Users_Change_Billing_ContactInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Change_Billing_ContactInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_change_billing_contact(inputs)
	return __es.users_change_billing_contact(inputs)
});
/**
* | output |
* | --- |
* | "Convert" |
*
* @param {Users_ConvertInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_convert = /** @type {((inputs?: Users_ConvertInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_ConvertInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_convert(inputs)
	return __es.users_convert(inputs)
});
/**
* | output |
* | --- |
* | "Convert to member" |
*
* @param {Users_Convert_To_MemberInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_convert_to_member = /** @type {((inputs?: Users_Convert_To_MemberInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Convert_To_MemberInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_convert_to_member(inputs)
	return __es.users_convert_to_member(inputs)
});
/**
* | output |
* | --- |
* | "Convert {user} to {role}" |
*
* @param {Users_Convert_User_To_RoleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_convert_user_to_role = /** @type {((inputs: Users_Convert_User_To_RoleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Convert_User_To_RoleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_convert_user_to_role(inputs)
	return __es.users_convert_user_to_role(inputs)
});
/**
* | output |
* | --- |
* | "Copy URL" |
*
* @param {Users_Copy_UrlInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_copy_url = /** @type {((inputs?: Users_Copy_UrlInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Copy_UrlInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_copy_url(inputs)
	return __es.users_copy_url(inputs)
});
/**
* | output |
* | --- |
* | "Create" |
*
* @param {Users_CreateInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_create = /** @type {((inputs?: Users_CreateInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_CreateInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_create(inputs)
	return __es.users_create(inputs)
});
/**
* | output |
* | --- |
* | "Email or group, separated by commas" |
*
* @param {Users_Email_Or_Group_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_email_or_group_placeholder = /** @type {((inputs?: Users_Email_Or_Group_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Email_Or_Group_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_email_or_group_placeholder(inputs)
	return __es.users_email_or_group_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Add emails, separated by commas" |
*
* @param {Users_Email_PlaceholderInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_email_placeholder = /** @type {((inputs?: Users_Email_PlaceholderInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Email_PlaceholderInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_email_placeholder(inputs)
	return __es.users_email_placeholder(inputs)
});
/**
* | output |
* | --- |
* | "Error" |
*
* @param {Users_ErrorInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_error = /** @type {((inputs?: Users_ErrorInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_ErrorInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_error(inputs)
	return __es.users_error(inputs)
});
/**
* | output |
* | --- |
* | "Error loading organization members:" |
*
* @param {Users_Error_Loading_MembersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_error_loading_members = /** @type {((inputs?: Users_Error_Loading_MembersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Error_Loading_MembersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_error_loading_members(inputs)
	return __es.users_error_loading_members(inputs)
});
/**
* | output |
* | --- |
* | "Error removing user from organization" |
*
* @param {Users_Error_RemovingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_error_removing = /** @type {((inputs?: Users_Error_RemovingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Error_RemovingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_error_removing(inputs)
	return __es.users_error_removing(inputs)
});
/**
* | output |
* | --- |
* | "Error removing user" |
*
* @param {Users_Error_Removing_UserInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_error_removing_user = /** @type {((inputs?: Users_Error_Removing_UserInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Error_Removing_UserInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_error_removing_user(inputs)
	return __es.users_error_removing_user(inputs)
});
/**
* | output |
* | --- |
* | "Error updating user role" |
*
* @param {Users_Error_Updating_RoleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_error_updating_role = /** @type {((inputs?: Users_Error_Updating_RoleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Error_Updating_RoleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_error_updating_role(inputs)
	return __es.users_error_updating_role(inputs)
});
/**
* | output |
* | --- |
* | "Error upgrading user role" |
*
* @param {Users_Error_Upgrading_RoleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_error_upgrading_role = /** @type {((inputs?: Users_Error_Upgrading_RoleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Error_Upgrading_RoleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_error_upgrading_role(inputs)
	return __es.users_error_upgrading_role(inputs)
});
/**
* | output |
* | --- |
* | "Everyone at {organization}" |
*
* @param {Users_Everyone_At_OrgInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_everyone_at_org = /** @type {((inputs: Users_Everyone_At_OrgInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Everyone_At_OrgInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_everyone_at_org(inputs)
	return __es.users_everyone_at_org(inputs)
});
/**
* | output |
* | --- |
* | "Failed to add groups: {groups}" |
*
* @param {Users_Failed_Add_GroupsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_failed_add_groups = /** @type {((inputs: Users_Failed_Add_GroupsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Failed_Add_GroupsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_failed_add_groups(inputs)
	return __es.users_failed_add_groups(inputs)
});
/**
* | output |
* | --- |
* | "Failed to invite {emails}" |
*
* @param {Users_Failed_InviteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_failed_invite = /** @type {((inputs: Users_Failed_InviteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Failed_InviteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_failed_invite(inputs)
	return __es.users_failed_invite(inputs)
});
/**
* | output |
* | --- |
* | "Failed to invite users: {emails}" |
*
* @param {Users_Failed_Invite_UsersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_failed_invite_users = /** @type {((inputs: Users_Failed_Invite_UsersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Failed_Invite_UsersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_failed_invite_users(inputs)
	return __es.users_failed_invite_users(inputs)
});
/**
* | output |
* | --- |
* | "Failed to load projects" |
*
* @param {Users_Failed_Load_ProjectsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_failed_load_projects = /** @type {((inputs?: Users_Failed_Load_ProjectsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Failed_Load_ProjectsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_failed_load_projects(inputs)
	return __es.users_failed_load_projects(inputs)
});
/**
* | output |
* | --- |
* | "Failed to load users" |
*
* @param {Users_Failed_Load_UsersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_failed_load_users = /** @type {((inputs?: Users_Failed_Load_UsersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Failed_Load_UsersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_failed_load_users(inputs)
	return __es.users_failed_load_users(inputs)
});
/**
* | output |
* | --- |
* | "Admins" |
*
* @param {Users_Filter_AdminsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_filter_admins = /** @type {((inputs?: Users_Filter_AdminsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Filter_AdminsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_filter_admins(inputs)
	return __es.users_filter_admins(inputs)
});
/**
* | output |
* | --- |
* | "All" |
*
* @param {Users_Filter_AllInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_filter_all = /** @type {((inputs?: Users_Filter_AllInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Filter_AllInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_filter_all(inputs)
	return __es.users_filter_all(inputs)
});
/**
* | output |
* | --- |
* | "All Roles" |
*
* @param {Users_Filter_All_RolesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_filter_all_roles = /** @type {((inputs?: Users_Filter_All_RolesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Filter_All_RolesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_filter_all_roles(inputs)
	return __es.users_filter_all_roles(inputs)
});
/**
* | output |
* | --- |
* | "All users" |
*
* @param {Users_Filter_All_UsersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_filter_all_users = /** @type {((inputs?: Users_Filter_All_UsersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Filter_All_UsersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_filter_all_users(inputs)
	return __es.users_filter_all_users(inputs)
});
/**
* | output |
* | --- |
* | "Editors" |
*
* @param {Users_Filter_EditorsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_filter_editors = /** @type {((inputs?: Users_Filter_EditorsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Filter_EditorsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_filter_editors(inputs)
	return __es.users_filter_editors(inputs)
});
/**
* | output |
* | --- |
* | "Members" |
*
* @param {Users_Filter_MembersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_filter_members = /** @type {((inputs?: Users_Filter_MembersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Filter_MembersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_filter_members(inputs)
	return __es.users_filter_members(inputs)
});
/**
* | output |
* | --- |
* | "Pending invites" |
*
* @param {Users_Filter_Pending_InvitesInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_filter_pending_invites = /** @type {((inputs?: Users_Filter_Pending_InvitesInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Filter_Pending_InvitesInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_filter_pending_invites(inputs)
	return __es.users_filter_pending_invites(inputs)
});
/**
* | output |
* | --- |
* | "Viewers" |
*
* @param {Users_Filter_ViewersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_filter_viewers = /** @type {((inputs?: Users_Filter_ViewersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Filter_ViewersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_filter_viewers(inputs)
	return __es.users_filter_viewers(inputs)
});
/**
* | output |
* | --- |
* | "Name" |
*
* @param {Users_Form_NameInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_form_name = /** @type {((inputs?: Users_Form_NameInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Form_NameInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_form_name(inputs)
	return __es.users_form_name(inputs)
});
/**
* | output |
* | --- |
* | "Untitled" |
*
* @param {Users_Form_UntitledInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_form_untitled = /** @type {((inputs?: Users_Form_UntitledInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Form_UntitledInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_form_untitled(inputs)
	return __es.users_form_untitled(inputs)
});
/**
* | output |
* | --- |
* | "Users" |
*
* @param {Users_Form_UsersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_form_users = /** @type {((inputs?: Users_Form_UsersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Form_UsersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_form_users(inputs)
	return __es.users_form_users(inputs)
});
/**
* | output |
* | --- |
* | "General Access" |
*
* @param {Users_General_AccessInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_general_access = /** @type {((inputs?: Users_General_AccessInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_General_AccessInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_general_access(inputs)
	return __es.users_general_access(inputs)
});
/**
* | output |
* | --- |
* | "{count} Groups" |
*
* @param {Users_Group_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_group_count = /** @type {((inputs: Users_Group_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Group_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_group_count(inputs)
	return __es.users_group_count(inputs)
});
/**
* | output |
* | --- |
* | "Guest upgraded to member and assigned {role} role" |
*
* @param {Users_Guest_UpgradedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_guest_upgraded = /** @type {((inputs: Users_Guest_UpgradedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Guest_UpgradedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_guest_upgraded(inputs)
	return __es.users_guest_upgraded(inputs)
});
/**
* | output |
* | --- |
* | "Guest upgraded to {role}" |
*
* @param {Users_Guest_Upgraded_ToInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_guest_upgraded_to = /** @type {((inputs: Users_Guest_Upgraded_ToInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Guest_Upgraded_ToInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_guest_upgraded_to(inputs)
	return __es.users_guest_upgraded_to(inputs)
});
/**
* | output |
* | --- |
* | "Invited {count} guests as {role}" |
*
* @param {Users_Guests_Invited_SuccessInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_guests_invited_success = /** @type {((inputs: Users_Guests_Invited_SuccessInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Guests_Invited_SuccessInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_guests_invited_success(inputs)
	return __es.users_guests_invited_success(inputs)
});
/**
* | output |
* | --- |
* | "Invite" |
*
* @param {Users_InviteInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_invite = /** @type {((inputs?: Users_InviteInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_InviteInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_invite(inputs)
	return __es.users_invite(inputs)
});
/**
* | output |
* | --- |
* | "Invite only" |
*
* @param {Users_Invite_OnlyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_invite_only = /** @type {((inputs?: Users_Invite_OnlyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Invite_OnlyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_invite_only(inputs)
	return __es.users_invite_only(inputs)
});
/**
* | output |
* | --- |
* | "Invited {count} people" |
*
* @param {Users_Invited_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_invited_count = /** @type {((inputs: Users_Invited_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Invited_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_invited_count(inputs)
	return __es.users_invited_count(inputs)
});
/**
* | output |
* | --- |
* | "Successfully invited {count} people as {role}" |
*
* @param {Users_Invited_SuccessInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_invited_success = /** @type {((inputs: Users_Invited_SuccessInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Invited_SuccessInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_invited_success(inputs)
	return __es.users_invited_success(inputs)
});
/**
* | output |
* | --- |
* | "Learn more about sharing" |
*
* @param {Users_Learn_More_SharingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_learn_more_sharing = /** @type {((inputs?: Users_Learn_More_SharingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Learn_More_SharingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_learn_more_sharing(inputs)
	return __es.users_learn_more_sharing(inputs)
});
/**
* | output |
* | --- |
* | "Loading..." |
*
* @param {Users_LoadingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_loading = /** @type {((inputs?: Users_LoadingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_LoadingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_loading(inputs)
	return __es.users_loading(inputs)
});
/**
* | output |
* | --- |
* | "Loading more..." |
*
* @param {Users_Loading_MoreInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_loading_more = /** @type {((inputs?: Users_Loading_MoreInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Loading_MoreInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_loading_more(inputs)
	return __es.users_loading_more(inputs)
});
/**
* | output |
* | --- |
* | "{count} members" |
*
* @param {Users_Member_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_member_count = /** @type {((inputs: Users_Member_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Member_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_member_count(inputs)
	return __es.users_member_count(inputs)
});
/**
* | output |
* | --- |
* | "No groups" |
*
* @param {Users_No_GroupsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_no_groups = /** @type {((inputs?: Users_No_GroupsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_No_GroupsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_no_groups(inputs)
	return __es.users_no_groups(inputs)
});
/**
* | output |
* | --- |
* | "No projects" |
*
* @param {Users_No_ProjectsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_no_projects = /** @type {((inputs?: Users_No_ProjectsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_No_ProjectsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_no_projects(inputs)
	return __es.users_no_projects(inputs)
});
/**
* | output |
* | --- |
* | "No users" |
*
* @param {Users_No_UsersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_no_users = /** @type {((inputs?: Users_No_UsersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_No_UsersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_no_users(inputs)
	return __es.users_no_users(inputs)
});
/**
* | output |
* | --- |
* | "Only admins and invited users can access" |
*
* @param {Users_Only_Admins_AccessInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_only_admins_access = /** @type {((inputs?: Users_Only_Admins_AccessInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Only_Admins_AccessInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_only_admins_access(inputs)
	return __es.users_only_admins_access(inputs)
});
/**
* | output |
* | --- |
* | "Org members can access" |
*
* @param {Users_Org_Members_AccessInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_org_members_access = /** @type {((inputs?: Users_Org_Members_AccessInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Org_Members_AccessInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_org_members_access(inputs)
	return __es.users_org_members_access(inputs)
});
/**
* | output |
* | --- |
* | "Manage users" |
*
* @param {Users_Page_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_page_title = /** @type {((inputs?: Users_Page_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Page_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_page_title(inputs)
	return __es.users_page_title(inputs)
});
/**
* | output |
* | --- |
* | "Pending invitation" |
*
* @param {Users_Pending_InvitationInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_pending_invitation = /** @type {((inputs?: Users_Pending_InvitationInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Pending_InvitationInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_pending_invitation(inputs)
	return __es.users_pending_invitation(inputs)
});
/**
* | output |
* | --- |
* | "Project access" |
*
* @param {Users_Project_AccessInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_project_access = /** @type {((inputs?: Users_Project_AccessInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Project_AccessInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_project_access(inputs)
	return __es.users_project_access(inputs)
});
/**
* | output |
* | --- |
* | "{count} Projects" |
*
* @param {Users_Project_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_project_count = /** @type {((inputs: Users_Project_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Project_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_project_count(inputs)
	return __es.users_project_count(inputs)
});
/**
* | output |
* | --- |
* | "{count} Projects" |
*
* @param {Users_Projects_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_projects_count = /** @type {((inputs: Users_Projects_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Projects_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_projects_count(inputs)
	return __es.users_projects_count(inputs)
});
/**
* | output |
* | --- |
* | "Remove" |
*
* @param {Users_RemoveInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_remove = /** @type {((inputs?: Users_RemoveInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_RemoveInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_remove(inputs)
	return __es.users_remove(inputs)
});
/**
* | output |
* | --- |
* | "This user will no longer be able to access the organization." |
*
* @param {Users_Remove_Confirm_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_remove_confirm_desc = /** @type {((inputs?: Users_Remove_Confirm_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Remove_Confirm_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_remove_confirm_desc(inputs)
	return __es.users_remove_confirm_desc(inputs)
});
/**
* | output |
* | --- |
* | "Remove user from organization?" |
*
* @param {Users_Remove_Confirm_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_remove_confirm_title = /** @type {((inputs?: Users_Remove_Confirm_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Remove_Confirm_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_remove_confirm_title(inputs)
	return __es.users_remove_confirm_title(inputs)
});
/**
* | output |
* | --- |
* | "User removed" |
*
* @param {Users_RemovedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_removed = /** @type {((inputs?: Users_RemovedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_RemovedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_removed(inputs)
	return __es.users_removed(inputs)
});
/**
* | output |
* | --- |
* | "User removed from organization" |
*
* @param {Users_Removed_From_OrgInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_removed_from_org = /** @type {((inputs?: Users_Removed_From_OrgInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Removed_From_OrgInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_removed_from_org(inputs)
	return __es.users_removed_from_org(inputs)
});
/**
* | output |
* | --- |
* | "Retry" |
*
* @param {Users_RetryInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_retry = /** @type {((inputs?: Users_RetryInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_RetryInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_retry(inputs)
	return __es.users_retry(inputs)
});
/**
* | output |
* | --- |
* | "User role updated" |
*
* @param {Users_Role_UpdatedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_role_updated = /** @type {((inputs?: Users_Role_UpdatedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Role_UpdatedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_role_updated(inputs)
	return __es.users_role_updated(inputs)
});
/**
* | output |
* | --- |
* | "Save" |
*
* @param {Users_SaveInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_save = /** @type {((inputs?: Users_SaveInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_SaveInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_save(inputs)
	return __es.users_save(inputs)
});
/**
* | output |
* | --- |
* | "Searching..." |
*
* @param {Users_SearchingInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_searching = /** @type {((inputs?: Users_SearchingInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_SearchingInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_searching(inputs)
	return __es.users_searching(inputs)
});
/**
* | output |
* | --- |
* | "GROUPS" |
*
* @param {Users_Section_GroupsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_section_groups = /** @type {((inputs?: Users_Section_GroupsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Section_GroupsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_section_groups(inputs)
	return __es.users_section_groups(inputs)
});
/**
* | output |
* | --- |
* | "GUESTS" |
*
* @param {Users_Section_GuestsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_section_guests = /** @type {((inputs?: Users_Section_GuestsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Section_GuestsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_section_guests(inputs)
	return __es.users_section_guests(inputs)
});
/**
* | output |
* | --- |
* | "MEMBERS" |
*
* @param {Users_Section_MembersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_section_members = /** @type {((inputs?: Users_Section_MembersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Section_MembersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_section_members(inputs)
	return __es.users_section_members(inputs)
});
/**
* | output |
* | --- |
* | "Select projects" |
*
* @param {Users_Select_ProjectsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_select_projects = /** @type {((inputs?: Users_Select_ProjectsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Select_ProjectsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_select_projects(inputs)
	return __es.users_select_projects(inputs)
});
/**
* | output |
* | --- |
* | "Groups ({count})" |
*
* @param {Users_Tab_GroupsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_tab_groups = /** @type {((inputs: Users_Tab_GroupsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Tab_GroupsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_tab_groups(inputs)
	return __es.users_tab_groups(inputs)
});
/**
* | output |
* | --- |
* | "Guests ({count})" |
*
* @param {Users_Tab_GuestsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_tab_guests = /** @type {((inputs: Users_Tab_GuestsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Tab_GuestsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_tab_guests(inputs)
	return __es.users_tab_guests(inputs)
});
/**
* | output |
* | --- |
* | "Members ({count})" |
*
* @param {Users_Tab_MembersInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_tab_members = /** @type {((inputs: Users_Tab_MembersInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Tab_MembersInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_tab_members(inputs)
	return __es.users_tab_members(inputs)
});
/**
* | output |
* | --- |
* | "No users found" |
*
* @param {Users_Table_EmptyInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_table_empty = /** @type {((inputs?: Users_Table_EmptyInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Table_EmptyInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_table_empty(inputs)
	return __es.users_table_empty(inputs)
});
/**
* | output |
* | --- |
* | "Groups" |
*
* @param {Users_Table_Header_GroupsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_table_header_groups = /** @type {((inputs?: Users_Table_Header_GroupsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Table_Header_GroupsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_table_header_groups(inputs)
	return __es.users_table_header_groups(inputs)
});
/**
* | output |
* | --- |
* | "Organization Role" |
*
* @param {Users_Table_Header_Org_RoleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_table_header_org_role = /** @type {((inputs?: Users_Table_Header_Org_RoleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Table_Header_Org_RoleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_table_header_org_role(inputs)
	return __es.users_table_header_org_role(inputs)
});
/**
* | output |
* | --- |
* | "Projects" |
*
* @param {Users_Table_Header_ProjectsInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_table_header_projects = /** @type {((inputs?: Users_Table_Header_ProjectsInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Table_Header_ProjectsInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_table_header_projects(inputs)
	return __es.users_table_header_projects(inputs)
});
/**
* | output |
* | --- |
* | "User" |
*
* @param {Users_Table_Header_UserInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_table_header_user = /** @type {((inputs?: Users_Table_Header_UserInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Table_Header_UserInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_table_header_user(inputs)
	return __es.users_table_header_user(inputs)
});
/**
* | output |
* | --- |
* | "Upgrading a guest to {role} will grant this user access to all open projects in the organization. Would you like to upgrade this guest user to {role}?" |
*
* @param {Users_Upgrade_Confirm_DescInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_upgrade_confirm_desc = /** @type {((inputs: Users_Upgrade_Confirm_DescInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Upgrade_Confirm_DescInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_upgrade_confirm_desc(inputs)
	return __es.users_upgrade_confirm_desc(inputs)
});
/**
* | output |
* | --- |
* | "Upgrade guest to {role}?" |
*
* @param {Users_Upgrade_Confirm_TitleInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_upgrade_confirm_title = /** @type {((inputs: Users_Upgrade_Confirm_TitleInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Upgrade_Confirm_TitleInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_upgrade_confirm_title(inputs)
	return __es.users_upgrade_confirm_title(inputs)
});
/**
* | output |
* | --- |
* | "URL copied" |
*
* @param {Users_Url_CopiedInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_url_copied = /** @type {((inputs?: Users_Url_CopiedInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Url_CopiedInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_url_copied(inputs)
	return __es.users_url_copied(inputs)
});
/**
* | output |
* | --- |
* | "{count} Users" |
*
* @param {Users_User_CountInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_user_count = /** @type {((inputs: Users_User_CountInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_User_CountInputs, { locale?: "en" | "es" }, {}>} */ ((inputs, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_user_count(inputs)
	return __es.users_user_count(inputs)
});
/**
* | output |
* | --- |
* | "Yes, remove" |
*
* @param {Users_Yes_RemoveInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_yes_remove = /** @type {((inputs?: Users_Yes_RemoveInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Yes_RemoveInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_yes_remove(inputs)
	return __es.users_yes_remove(inputs)
});
/**
* | output |
* | --- |
* | "Yes, upgrade" |
*
* @param {Users_Yes_UpgradeInputs} inputs
* @param {{ locale?: "en" | "es" }} options
* @returns {LocalizedString}
*/
export const users_yes_upgrade = /** @type {((inputs?: Users_Yes_UpgradeInputs, options?: { locale?: "en" | "es" }) => LocalizedString) & import('../runtime.js').MessageMetadata<Users_Yes_UpgradeInputs, { locale?: "en" | "es" }, {}>} */ ((inputs = {}, options = {}) => {
	const locale = experimentalStaticLocale ?? options.locale ?? getLocale()
	if (locale === "en") return __en.users_yes_upgrade(inputs)
	return __es.users_yes_upgrade(inputs)
});