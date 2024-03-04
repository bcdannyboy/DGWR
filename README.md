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

The main file (`main.go`) in the DGWR system provides a detailed example of how to define risk events, including their probabilities, impacts, and dependencies. The example focuses on modeling a ransomware event tree, demonstrating the system's capability to handle complex, interconnected risk scenarios.

### Output Files

- **`probabilities.json`**: Contains the probability outputs for each scenario, derived from running the simulation 100,000 times. This file offers insight into the likelihood of various risk events occurring within the modeled scenario.
- **`impacts.json`**: Stores the impact values associated with each risk event, also based on 100,000 simulation runs. These values help quantify the potential consequences of risk events, aiding in risk assessment and mitigation planning.
