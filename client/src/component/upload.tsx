import React, { useState } from "react";
import { Button, Paper } from "@mui/material";
import { useQuery } from "react-query";
import { uploadFile } from "../api/getListResultScraping";

const FileUpload = () => {
  const [selectedFile, setSelectedFile] = useState<File | null>();

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files !== null) {
      const file = event.target.files[0];
      setSelectedFile(file);
    }
  };

  const HandleUpload = () => {
    const formData = new FormData();
    if (selectedFile) {
      const blob = new Blob([selectedFile]);
      formData.append("file", blob, selectedFile.name);
    
    }
    console.log(selectedFile);
    uploadFile(formData);

    setSelectedFile(null);
  };

  return (
    <Paper elevation={3} style={{ padding: "16px" }}>
      <input
        type="file"
        accept=".csv"
        onChange={handleFileChange}
        style={{ display: "none" }}
        id="file-upload"
        name="file"
      />
      <label htmlFor="file-upload">
        <Button variant="contained" component="span">
          Upload CSV
        </Button>
      </label>
      <Button
        variant="contained"
        color="primary"
        onClick={HandleUpload}
        disabled={!selectedFile}
        style={{ marginLeft: "16px" }}
      >
        Upload
      </Button>
      {selectedFile && <p>Selected file: {selectedFile.name}</p>}
    </Paper>
  );
};

export default FileUpload;
