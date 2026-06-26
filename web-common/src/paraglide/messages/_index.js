/* eslint-disable */
import { getLocale, experimentalStaticLocale } from "../runtime.js"

/** @typedef {import('../runtime.js').LocalizedString} LocalizedString */
/** @typedef {{}} Bignumber_Copy_ValueInputs */
/** @typedef {{}} Bignumber_Shift_ClickInputs */
/** @typedef {{}} Calendar_ApplyInputs */
/** @typedef {{ capLabel: NonNullable<unknown> }} Calendar_Range_Exceeds_LimitInputs */
/** @typedef {{}} Canvas_Add_WidgetInputs */
/** @typedef {{}} Canvas_Ai_Generating_ChartInputs */
/** @typedef {{}} Canvas_Ai_HintInputs */
/** @typedef {{}} Canvas_Ai_Is_EditingInputs */
/** @typedef {{}} Canvas_Ai_Manual_HintInputs */
/** @typedef {{}} Canvas_Ai_Write_ManuallyInputs */
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
/** @typedef {{}} Canvas_No_Valid_Metrics_ViewInputs */
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
/** @typedef {{}} Canvas_Sparkline_LabelInputs */
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
/** @typedef {{}} Chat_Happy_To_ExploreInputs */
/** @typedef {{}} Chat_Placeholder_AnalystInputs */
/** @typedef {{}} Common_ApplyInputs */
/** @typedef {{}} Common_CancelInputs */
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
/** @typedef {{}} Explore_Go_To_DashboardInputs */
/** @typedef {{}} Explore_Go_To_ExploreInputs */
/** @typedef {{ name: NonNullable<unknown> }} Explore_Go_To_NamedInputs */
/** @typedef {{}} Explore_Unable_To_OpenInputs */
/** @typedef {{}} Filter_Enter_Search_TermInputs */
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
/** @typedef {{}} Filter_Mode_ContainsInputs */
/** @typedef {{}} Filter_Mode_Contains_DescriptionInputs */
/** @typedef {{}} Filter_Mode_In_ListInputs */
/** @typedef {{}} Filter_Mode_In_List_DescriptionInputs */
/** @typedef {{}} Filter_Mode_SelectInputs */
/** @typedef {{}} Filter_Mode_Select_DescriptionInputs */
/** @typedef {{}} Filter_Paste_List_HintInputs */
/** @typedef {{}} Footer_Report_IssueInputs */
/** @typedef {{}} Footer_Rill_DeveloperInputs */
/** @typedef {{}} Footer_Shortcut_ClickInputs */
/** @typedef {{}} Footer_Unknown_VersionInputs */
/** @typedef {{}} Footer_VersionInputs */
/** @typedef {{}} Footer_View_DocumentationInputs */
/** @typedef {{}} Language_EnInputs */
/** @typedef {{}} Language_EsInputs */
/** @typedef {{}} Language_Switcher_LabelInputs */
/** @typedef {{}} Layout_Inspector_Panel_AriaInputs */
/** @typedef {{}} Leaderboard_Copy_ValueInputs */
/** @typedef {{}} Leaderboard_Expand_TableInputs */
/** @typedef {{}} Leaderboard_Expand_TooltipInputs */
/** @typedef {{}} Leaderboard_No_Available_ValuesInputs */
/** @typedef {{}} Leaderboard_Shift_ClickInputs */
/** @typedef {{}} Measure_Filter_ApplyInputs */
/** @typedef {{}} Measure_Filter_By_DimensionInputs */
/** @typedef {{}} Measure_Filter_Enter_NumberInputs */
/** @typedef {{}} Measure_Filter_Higher_ValueInputs */
/** @typedef {{}} Measure_Filter_Lower_ValueInputs */
/** @typedef {{}} Measure_Filter_Select_DimensionInputs */
/** @typedef {{}} Measure_Filter_ThresholdInputs */
/** @typedef {{}} Nav_Close_SidebarInputs */
/** @typedef {{}} Nav_Data_ExplorerInputs */
/** @typedef {{}} Nav_Show_SidebarInputs */
/** @typedef {{ grain: NonNullable<unknown> }} Pivot_Time_Dimension_HeaderInputs */
/** @typedef {{}} Pivot_Time_PrefixInputs */
/** @typedef {{ duration: NonNullable<unknown> }} Time_AgoInputs */
/** @typedef {{}} Time_All_TimeInputs */
/** @typedef {{}} Time_ComparingInputs */
/** @typedef {{}} Time_CustomInputs */
/** @typedef {{}} Time_Custom_RangeInputs */
/** @typedef {{}} Time_Enter_Time_RangeInputs */
/** @typedef {{ duration: NonNullable<unknown> }} Time_From_NowInputs */
/** @typedef {{}} Time_Grain_ByInputs */
/** @typedef {{}} Time_Grain_CompleteInputs */
/** @typedef {{}} Time_Grain_TimeInputs */
/** @typedef {{ duration: NonNullable<unknown> }} Time_Last_DurationInputs */
/** @typedef {{}} Time_Month_To_DateInputs */
/** @typedef {{}} Time_No_Comparison_PeriodInputs */
/** @typedef {{}} Time_Previous_MonthInputs */
/** @typedef {{}} Time_Previous_QuarterInputs */
/** @typedef {{}} Time_Previous_WeekInputs */
/** @typedef {{}} Time_Previous_YearInputs */
/** @typedef {{}} Time_Quarter_To_DateInputs */
/** @typedef {{}} Time_Ref_CompleteInputs */
/** @typedef {{}} Time_Ref_Complete_DataInputs */
/** @typedef {{}} Time_Ref_CurrentInputs */
/** @typedef {{}} Time_Ref_LatestInputs */
/** @typedef {{}} Time_Ref_NowInputs */
/** @typedef {{}} Time_TodayInputs */
/** @typedef {{}} Time_Unable_To_ParseInputs */
/** @typedef {{}} Time_Week_To_DateInputs */
/** @typedef {{}} Time_Year_To_DateInputs */
/** @typedef {{}} Time_YesterdayInputs */
import * as __en from "./en.js"
import * as __es from "./es.js"
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