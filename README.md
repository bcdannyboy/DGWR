# *D*ont *G*amble *W*ith *R*isk

**Don't Gamble With Risk (DGWR)** is a Monte Carlo simulation system designed for complex quantitative risk modeling, focusing on risks, costs, and benefits. Inspired by the FAIR (Factor Analysis of Information Risk) model and leveraging the principles of Monte Carlo simulations, DGWR offers the basis for a robust framework for analyzing and quantifying risk probabilities and impacts. This system could be particularly useful for organizations looking to make informed decisions based on quantitative risk analysis.

## Key Features

- **Quantitative Risk Analysis**: DGWR provides a structured approach to quantifying risks, enabling an understanding of potential impacts in numerical terms.
- **Monte Carlo Simulations**: By running simulations multiple times (e.g., 100,000 iterations), DGWR offers a comprehensive view of possible outcomes, helping to identify the probability of various events.
- **Flexible Event Modeling**: The system supports complex event trees, allowing for the modeling of events with varying timeframes, probabilities, impacts, and dependencies.

## Benefits of Quantitative Risk Analysis

Quantitative risk analysis offers several advantages over qualitative analysis, including:

- **Objective Decision Making**: By quantifying risks, organizations can make decisions based on data, reducing bias and subjectivity.
- **Prioritization of Risks**: Quantitative analysis helps in identifying and prioritizing risks based on their potential impact, enabling more effective risk management strategies.
- **Resource Allocation**: Understanding the numerical impact of risks allows for better allocation of resources to mitigate high-priority risks.

## Current Implementation and Future Directions

Currently, DGWR utilizes the PERT (Program Evaluation and Review Technique) distribution for modeling the probabilities and impacts of each risk event. This choice was made to balance the need for accuracy with computational efficiency. However, the system's design allows for the integration of other distributions in the future, enhancing its flexibility and applicability to a wide range of scenarios.

## System Usage

There are 2 main files which can be used as implementation examples.

### Ransomware Event Tree Example
`main_ransomware.go` is an example of a complex phishing -> major ransomware event event tree. This file can be used as a reference for implementing event trees and using the DGWR system to run simulations and analyze the results. This file also contains detailed code comments for the qualification and explanation of input variables.

This example also comes with 2 pre-generated output files:  `probabilities_ransomware.json` and `impacts_ransomware.json` which contain the probabilities and impacts of the ransomware event tree.

### Code Vulnerability Event Tree Example
`main_codevuln.go` contains a much more simplistic example of a code vulnerability event tree. This file can be used as a reference for implementing event trees and using the DGWR system to run simulations and analyze the results. 

This example comes with a single pre-generated output file, `probabilities_vulnerability.json` which contains the probabilities of the code vulnerability event tree nodes.