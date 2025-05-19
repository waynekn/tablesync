// Represents data required to create a spreadsheet
export type SpreadsheetInit = {
  title: string;
  description: string;
  deadline: string;
  colTitles: string[];
};

export type SpreadsheetInitErr = Omit<Partial<SpreadsheetInit>, "colTitles"> & {
  // error key is a non field error
  error?: string;
  colTitles?: string;
};
