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

// Represents a spreadsheet retrieved from the api
export type Spreadsheet = {
  id: string;
  title: string;
  description: string;
  owner: string;
  createdAt: string;
  updatedAt: string;
  data: string[][];
  deadline: string;
};
