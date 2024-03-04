package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bcdannyboy/dgws/risk"
	"github.com/bcdannyboy/dgws/risk/analysis"
	"github.com/bcdannyboy/dgws/risk/utils"
)

// This is an example of a ransomware event tree stemming from a phishing email
// The goal of this example is to show how the risk package can be used to model the probabilities and impacts of complex event trees
// Each event in the tree has their own expected time frame, probability, impacts, and dependencies.
// All timeframes are scaled up to a yearly basis.
//
// In this example, Beta distributions and Latin Hypercube Sampling are used to model the probabilities and impacts of each event.
//
// In the output, all probabilities and impacts are represented in a scaled yearly frequency basis, regardless of the original timeframe.
//
// 								   +--------------------------------+
// 								   |                                |
// 								   |    Attempted Phishing Emails   |
// 								   |                                |
// 								   +----------------|---------------+
// 										    |
// 				  		   Anti-Phish Filter -------        |
// 								            \-------+
// 								                    |
// 								                    |       ------- Employee Reports Phish
// 										    +------/
// 										    |
// 										    |
// 								   +----------------|---------------+
// 								   |                                |
// 								   |   Employee's Account Details   |
// 								   |    Are Successfully Phished    |
// 								   |                                |
// 								   +----------------|---------------+
// 										    |
//                                 +-- Behavioural Controls Identify --------       |
//                                 |    Anomalous Account Activity           \------+
//                                 |                |     |                         |
//                                 |                |     |                         |
//                                 |                |     |        +----------------|---------------+
//                                 |                |     |        |                                |
//                                 |                |     |        |   Employee Accepts Malicious   |
//                                 |                |     |        |           Duo Prompt           |
//                                 |                |     |        |                                |
//                                 |                |     |        +----------------|---------------+
//                                 |                |     |                         |
//                                 |                |     |                         |       ------- Host Based Controls Identify
//                                 |                |     --------------------------+------/         Malicious Activity or Code
//                                 |                |                               |                     |  |
//                                 |                |                               |                     |  |
//                                 |                |              +----------------|---------------+     |  |
//                                 |                |              |                                |     |  |
//                                 |                |              |    Threat Actor Establishes    |     |  |
//                                 |                |              |         Code Execution         |     |  |
//                                 |                |              |                                |     |  |
//                                 |                |              +----------------|---------------+     |  |
//                                 |                +----------------------\        |                     |  |
//                                 |                                        -----\  |                     |  |
//                                 |        Network Controls Identify         ------+---------------------+  |
//                                 |        Command & Control Activity ------/      |                        |
//                                 |                         |                      |                        |
//                                 |                         |     +----------------|---------------+        |
//                                 |                         |     |                                |        |
//                                 |                         |     |      Threat Actor Delivers     |        |
//                                 |                         |     |        Ransomware Payload      |        |
//                                 |                         |     |                                |        |
//                                 |                         |     +----------------|---------------+        |
//                                 |                         |                      |                        |
//                                 |                         |                      |                        |
//                                 +-------------------------+----------------------+------------------------+
// 										    |
// 										    |
// 								   +----------------|---------------+
// 								   |                                |
// 								   |    Ransomware Propogates on    |
// 								   |    Network Without Detection   |
// 								   |                                |
// 								   +----------------|---------------+
// 								   		    |
// 								   		    |
// 								   		    |
// 								   +----------------|---------------+
// 								   |                                |
// 								   |     MAJOR RANSOMWARE EVENT     |
// 								   |                                |
// 								   +--------------------------------+

func main() {

	// --- ID GENERATION, SKIP AHEAD TO SEE THE EVENT DEFINITIONS ---
	EmployeeGetsPhishedID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for phishing event: %s", err.Error()))
	}

	PhishingID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for phishing event: %s", err.Error()))
	}

	AntiPhishFilterID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for anti-phishing filter event: %s", err.Error()))
	}

	EmployeeReportsPhishID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for employee reports phishing event: %s", err.Error()))
	}

	EmployeeAcceptsMaliciousDuoPushID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for employee accepts malicious Duo Push event: %s", err.Error()))
	}

	BehavioralControlsCatchAnomalousAccountBehaviorID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for behavioral controls catch anomalous account behavior event: %s", err.Error()))
	}

	ThreatActorEstablishesCodeExecutionCapabilitiesID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for threat actor establishes code execution capabilities event: %s", err.Error()))
	}

	HostBasedControlsCatchMaliciousActivityOrCodeID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for host-based controls catch malicious activity or code event: %s", err.Error()))
	}

	NetworkBasedControlsCatchMaliciousCommandAndControlTrafficID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for network-based controls catch malicious command and control traffic event: %s", err.Error()))
	}

	ThreatActorDeliversRansomwarePayloadID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for threat actor delivers ransomware payload event: %s", err.Error()))
	}

	RansomwarePropogatesThroughoutNetworkWithoutDetectionID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for ransomware propogates throughout network without detection event: %s", err.Error()))
	}

	MajorRansomwareEventID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for major ransomware event: %s", err.Error()))
	}

	impactPhishingEmailsDetectedID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for phishing emails detected impact: %s", err.Error()))
	}

	impactPhishingEmailsReportedID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for phishing emails reported impact: %s", err.Error()))
	}

	impactBehavioralAnomalyAlertsID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for behavioral anomaly alerts impact: %s", err.Error()))
	}

	impactHostControlAlertsID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for host control alerts impact: %s", err.Error()))
	}

	impactNetworkControlAlertsID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for network control alerts impact: %s", err.Error()))
	}

	impactCompromisedAccountsID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for compromised accounts impact: %s", err.Error()))
	}

	impactMaliciousDuoPromptsAcceptedID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for malicious Duo prompts accepted impact: %s", err.Error()))
	}

	impactPhishingEmailsReceivedID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for phishing emails received impact: %s", err.Error()))
	}

	impactThreatActorsWithAccessID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for threat actors with access impact: %s", err.Error()))
	}

	impactMalwareOnNetworkID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for malware on network impact: %s", err.Error()))
	}

	impactLateralMovementEventsID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for lateral movement events impact: %s", err.Error()))
	}

	impactRebuildingNetworkID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for rebuilding network impact: %s", err.Error()))
	}

	impactCustomerRenumerationID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for customer renumeration impact: %s", err.Error()))
	}

	impactCustomerLossID, err := utils.GenerateID()
	if err != nil {
		panic(fmt.Errorf("error generating ID for customer loss impact: %s", err.Error()))
	}
	// --- EVENT DEFINITIONS ---

	// Controls
	AntiPhishFilter := &risk.Event{
		ID:          AntiPhishFilterID,
		Name:        "Anti-Phishing Filter",
		Description: "An anti-phishing filter blocks a phishing email.",
		Probability: &risk.Probability{
			// We are 90% confident that the anti-phishing filter will block a phishing email at least 80% of the time, and at most 90% of the time
			// We expect phishing emails very commonly, so we're considering these predictions on a weekly frequency
			ExpectedFrequency: "weekly",
			Minimum:           0.8,
			MinimumConfidence: 0.9,
			Maximum:           0.9,
			MaximumConfidence: 0.9,
		},
		Impact: []*risk.Impact{
			{
				// The Threat Detection & Response group wants to track phishing emails detected as a metric.
				// For each filter event, the team expects to detect exactly 1 phishing email, so the minimum and maximum individual unit impacts are both 1
				// The team is very confident in their predictions, but want to account for any weirdness that might happen, so they've marked both their minimum and maximum predictions at 90% confidence
				ImpactID:                              impactPhishingEmailsDetectedID,
				Name:                                  "Phishing Emails Detected",
				Unit:                                  "Phishing Email Detected",
				PositiveImpact:                        true,
				Description:                           "The number of phishing emails detected by the anti-phishing filter.",
				ExpectedFrequency:                     "weekly",
				MinimumIndividualUnitImpact:           1,
				MinimumIndividualUnitImpactConfidence: 0.9,
				MaximumIndividualUnitImpact:           1,
				MaximumIndividualUnitImpactConfidence: 0.9,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.9,
				MaximumImpactEvents:                   1,
				MaximumImpactEventsConfidence:         0.9,
			},
		},
	}

	EmployeeReportsPhish := &risk.Event{
		ID:          EmployeeReportsPhishID,
		Name:        "Employee Reports Phishing",
		Description: "An employee reports a phishing email.",
		Probability: &risk.Probability{
			// Based on historical data, we find that employees report phishing emails at least 20% of the time, and at most 60% of the time
			// The team does phishing assessments once a quarter, so this data is on a quarterly frequency
			ExpectedFrequency: "quarterly",
			Minimum:           0.2,
			MinimumConfidence: 0.9,
			Maximum:           0.6,
			MaximumConfidence: 0.9,
		},
		Impact: []*risk.Impact{
			{
				// The Threat Detection & Response group wants to track phishing emails reported as a metric.
				// For each report event, the team expects to report exactly 1 phishing email, so the minimum and maximum individual unit impacts are both 1
				// The team is very confident in their predictions, but want to account for any weirdness that might happen, so they've marked both their minimum and maximum predictions at 90% confidence
				ImpactID:                              impactPhishingEmailsReportedID,
				Name:                                  "Phishing Emails Reported",
				Unit:                                  "Phishing Email Reported",
				PositiveImpact:                        true,
				Description:                           "The number of phishing emails reported by employees.",
				ExpectedFrequency:                     "quarterly",
				MinimumIndividualUnitImpact:           1,
				MinimumIndividualUnitImpactConfidence: 0.9,
				MaximumIndividualUnitImpact:           1,
				MaximumIndividualUnitImpactConfidence: 0.9,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.9,
				MaximumImpactEvents:                   1,
				MaximumImpactEventsConfidence:         0.9,
			},
		},
	}

	BehavioralControlsCatchAnomalousAccountBehavior := &risk.Event{
		ID:          BehavioralControlsCatchAnomalousAccountBehaviorID,
		Name:        "Behavioral Controls Catch Anomalous Account Behavior",
		Description: "Behavioral controls are triggered by anomalous account behavior.",
		Probability: &risk.Probability{
			// The Threat Detection and Response team just recently got a new behavioral control set, they're not sure how well it will work, but they're confident it will catch anomalous account behavior at least 50% of the time, and at most 80% of the time
			// The team is less confident in their maximum prediction, so they've marked it at 80% confidence but they're more confident in their minimum prediction, so they've marked it at 90% confidence
			// The team expects to get alerts about potentially anomalous account behavior throughout each month but they're unsure of the exact frequency within any given month, so they've marked it as a monthly frequency
			ExpectedFrequency: "monthly",
			Minimum:           0.5,
			MinimumConfidence: 0.9,
			Maximum:           0.8,
			MaximumConfidence: 0.8,
		},
		Impact: []*risk.Impact{
			{
				// The Threat Detection & Response group wants to track behavioral anomaly alerts as a metric.
				// For each event, the team knows that the behavioral controls will generate anywhere between 1 and 10 alerts, so they've marked their minimum and maximum impact events impacts at 1 and 10 respectively
				// The team is very confident in their predictions, but want to account for any weirdness that might happen, so they've marked both their minimum and maximum predictions at 90% confidence
				ImpactID:                              impactBehavioralAnomalyAlertsID,
				Name:                                  "Behavioral Anomaly Alerts",
				Unit:                                  "Behavioral Anomaly Alert",
				PositiveImpact:                        true,
				Description:                           "The number of alerts generated by the behavioral controls.",
				ExpectedFrequency:                     "monthly",
				MinimumIndividualUnitImpact:           1,
				MinimumIndividualUnitImpactConfidence: 0.9,
				MaximumIndividualUnitImpact:           1,
				MaximumIndividualUnitImpactConfidence: 0.9,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.9,
				MaximumImpactEvents:                   10,
				MaximumImpactEventsConfidence:         0.9,
			},
		},
	}

	HostBasedControlsCatchMaliciousActivityOrCode := &risk.Event{
		ID:          HostBasedControlsCatchMaliciousActivityOrCodeID,
		Name:        "Host-Based Controls Catch Malicious Activity or Code",
		Description: "Host-based controls catch malicious activity or code.",
		Probability: &risk.Probability{
			// The Threat Detection & Response team is confident that their host-based controls will catch malicious activity or code at least 60% of the time, and at most 90% of the time
			// The team is very confident in their predictions, as they've worked dilligently to test and implement these controls, so they've marked both their minimum and maximum predictions at 90% confidence
			// The team expects to get alerts about potentially malicious activity or code throughout each week, so they've marked it as a weekly frequency
			ExpectedFrequency: "weekly",
			Minimum:           0.6,
			MinimumConfidence: 0.9,
			Maximum:           0.90,
			MaximumConfidence: 0.9,
		},
		Impact: []*risk.Impact{
			{
				// The Threat Detection & Response group wants to track host control alerts as a metric.
				// For each event, the team knows that the host-based controls will generate between 1 and 10 alerts, so they've marked their minimum and maximum impact events impacts at 1 and 10 respectively
				// The team is very confident in their predictions, but want to account for any weirdness that might happen, so they've marked both their minimum and maximum predictions at 90% confidence
				ImpactID:                              impactHostControlAlertsID,
				Name:                                  "Host Control Alerts",
				Unit:                                  "Host Control Alert",
				PositiveImpact:                        true,
				Description:                           "The number of alerts generated by the host-based controls.",
				ExpectedFrequency:                     "weekly",
				MinimumIndividualUnitImpact:           1,
				MinimumIndividualUnitImpactConfidence: 0.9,
				MaximumIndividualUnitImpact:           1,
				MaximumIndividualUnitImpactConfidence: 0.9,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.9,
				MaximumImpactEvents:                   10,
				MaximumImpactEventsConfidence:         0.9,
			},
		},
	}

	NetworkBasedControlsCatchMaliciousCommandAndControlTraffic := &risk.Event{
		ID:          NetworkBasedControlsCatchMaliciousCommandAndControlTrafficID,
		Name:        "Network-Based Controls Catch Malicious Command and Control Traffic",
		Description: "Network-based controls catch malicious command and control traffic.",
		Probability: &risk.Probability{
			// The Threat Detection & Response team is less confident in their network-based controls, as they have multiple known gaps
			// The team is confident that their network-based controls will catch malicious command and control traffic at least 40% of the time, and at most 70% of the time
			// The team is not very confident in their maximum prediction, so they've marked it at 60% confidence but they're more confident in their minimum prediction, so they've marked it at 80% confidence
			// The team expects to get alerts about potentially malicious command and control traffic throughout each month but they're unsure of the exact frequency within any given month, so they've marked it as a monthly frequency
			ExpectedFrequency: "monthly",
			Minimum:           0.4,
			MinimumConfidence: 0.8,
			Maximum:           0.7,
			MaximumConfidence: 0.6,
		},
		Impact: []*risk.Impact{
			{
				// The Threat Detection & Response group wants to track network control alerts as a metric.
				// For each event, the team knows that the network-based controls will generate between 1 and 10 alerts, so they've marked their minimum and maximum impact events impacts at 1 and 10 respectively
				// The team is very confident in their predictions, but want to account for any weirdness that might happen, so they've marked both their minimum and maximum predictions at 90% confidence
				ImpactID:                              impactNetworkControlAlertsID,
				Name:                                  "Network Control Alerts",
				Unit:                                  "Network Control Alert",
				PositiveImpact:                        true,
				Description:                           "The number of alerts generated by the network-based controls.",
				ExpectedFrequency:                     "monthly",
				MinimumIndividualUnitImpact:           1,
				MinimumIndividualUnitImpactConfidence: 0.9,
				MaximumIndividualUnitImpact:           1,
				MaximumIndividualUnitImpactConfidence: 0.9,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.9,
				MaximumImpactEvents:                   10,
				MaximumImpactEventsConfidence:         0.9,
			},
		},
	}

	// Internal Employee Interaction Requirements
	EmployeeGetsPhished := &risk.Event{
		ID:          EmployeeGetsPhishedID,
		Name:        "Employee Falls for Phishing Email",
		Description: "An employee falls for a phishing email.",
		Probability: &risk.Probability{
			// Based on historical data, the team knows employees report phishing between 20% and 60% of the time
			// This means that employees fall for phishing emails between 40% and 80% of the time
			// since the team also wants to account for simply not reporting phishing emails, they've marked down their confidence to 70% for both their minimum and maximum predictions
			// The team does phishing assessments once a quarter, so this data is on a quarterly frequency
			ExpectedFrequency: "quarterly",
			Minimum:           0.2,
			MinimumConfidence: 0.7,
			Maximum:           0.6,
			MaximumConfidence: 0.7,
		},
		Impact: []*risk.Impact{
			{
				// The Threat Detection & Response group wants to estimate the number of compromised accounts as a result of successful phishing emails.
				// on average users have between 1 and 3 accounts when the team considers non-unique IDs or segmented admin accounts.
				// The team is very confident in their predictions, but want to account for any weirdness that might happen, so they've marked both their minimum and maximum predictions at 90% confidence
				// The team knows that users often re-use passwords so they think anywhere from 1 to 3 accounts can be compromised by a single phishing email
				// The team is a bit hopeful that password re-use isn't SO common so they've marked their maximum individual unit impact confidence at 70%
				ImpactID:                              impactCompromisedAccountsID,
				Name:                                  "Compromised Accounts",
				Unit:                                  "Compromised Account",
				PositiveImpact:                        false,
				Description:                           "The number of accounts compromised as a result of successful phishing emails.",
				ExpectedFrequency:                     "quarterly",
				MinimumIndividualUnitImpact:           1,
				MinimumIndividualUnitImpactConfidence: 0.9,
				MaximumIndividualUnitImpact:           3,
				MaximumIndividualUnitImpactConfidence: 0.7,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.9,
				MaximumImpactEvents:                   3,
				MaximumImpactEventsConfidence:         0.9,
			},
		},
		Dependencies: []*risk.Dependency{
			{
				// A phishing attempt has to happen for an employee to get phished
				DependsOnEventID: PhishingID,
				Happens:          true,
			},
			{
				// The anti-phishing filter has to not block the phishing email for the employee to get phished
				DependsOnEventID: AntiPhishFilterID,
				Happens:          false,
			},
		},
	}

	EmployeeAcceptsMaliciousDuoPush := &risk.Event{
		ID:          EmployeeAcceptsMaliciousDuoPushID,
		Name:        "Employee Accepts Malicious Duo Push",
		Description: "An employee accepts a malicious Duo Push notification.",
		Probability: &risk.Probability{
			// The Threat Detection & Response team has not directly observed an employee accepting a malicious Duo Push notification, but they have observed employees accepting legitimate Duo Push notifications without thinking twice about it or verifying the request
			// The team thinks that an employee will accept a malicious Duo Push notification at least 30% of the time, and at most 70% of the time
			// The team is really unsure about their maximum prediction, so they've marked it at a 50% (meaning it could be a coin flip if they're right or wrong), but they're a bit more confident in their minimum prediction, so they've marked it at 70% confidence
			ExpectedFrequency: "yearly",
			Minimum:           0.3,
			MinimumConfidence: 0.7,
			Maximum:           0.7,
			MaximumConfidence: 0.5,
		},
		Impact: []*risk.Impact{
			{
				// The Threat Detection & Response group wants to estimate the number of accounts compromised as a result of accepting a malicious Duo Push notification.
				// The team knows that despite user's having multiple accounts, duo push notifications are only ever used to authenticate a single account
				// The team is very confident in their predictions, but want to account for any weirdness that might happen, so they've marked both their minimum and maximum predictions at 90% confidence
				ImpactID:                              impactMaliciousDuoPromptsAcceptedID,
				Name:                                  "Malicious Duo Prompts Accepted",
				Unit:                                  "Malicious Duo Prompt Accepted",
				PositiveImpact:                        false,
				Description:                           "The number of accounts compromised as a result of accepting a malicious Duo Push notification.",
				ExpectedFrequency:                     "yearly",
				MinimumIndividualUnitImpact:           1,
				MinimumIndividualUnitImpactConfidence: 0.9,
				MaximumIndividualUnitImpact:           1,
				MaximumIndividualUnitImpactConfidence: 0.9,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.9,
				MaximumImpactEvents:                   1,
				MaximumImpactEventsConfidence:         0.9,
			},
		},
		Dependencies: []*risk.Dependency{
			{
				// The employee has to get phished for them to accept a malicious Duo Push notification
				DependsOnEventID: EmployeeGetsPhishedID,
				Happens:          true,
			},
			{
				// The behavioral controls have to not catch the anomalous account behavior for the employee to accept a malicious Duo Push notification
				DependsOnEventID: BehavioralControlsCatchAnomalousAccountBehaviorID,
				Happens:          false,
			},
		},
	}

	// Threat Actor Activities
	AttemptedPhishingEmail := &risk.Event{
		ID:          PhishingID,
		Name:        "Phishing Attempt",
		Description: "A threat actor sends a phishing email.",
		Probability: &risk.Probability{
			// The Threat Detection & Response team observes attempted phishing emails on a daily to every other day basis
			// The team is confident that they will observe attempted phishing emails at least 50% of the time, and at most 100% of the time (i.e. they either see some or they don't on any given day)
			ExpectedFrequency: "daily",
			Minimum:           0.5,
			MinimumConfidence: 0.9,
			Maximum:           1,
			MaximumConfidence: 0.9,
		},
		Impact: []*risk.Impact{
			{
				// The Threat Detection & Response group wants to estimate the number of phishing emails received as a result of attempted phishing emails.
				// The team knows that on average, users receive between 1 and 5 phishing emails a day, so they've marked their minimum and maximum impact events at 1 and 5 respectively
				// The team is very confident in their predictions, but want to account for any weirdness that might happen, so they've marked both their minimum and maximum predictions at 90% confidence
				ImpactID:                              impactPhishingEmailsReceivedID,
				Name:                                  "Phishing Emails Received",
				Unit:                                  "Phishing Email Received",
				PositiveImpact:                        false,
				Description:                           "The number of phishing emails received as a result of attempted phishing emails.",
				ExpectedFrequency:                     "daily",
				MinimumIndividualUnitImpact:           1,
				MinimumIndividualUnitImpactConfidence: 0.9,
				MaximumIndividualUnitImpact:           5,
				MaximumIndividualUnitImpactConfidence: 0.9,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.9,
				MaximumImpactEvents:                   5,
				MaximumImpactEventsConfidence:         0.9,
			},
		},
	}

	ThreatActorEstablishesCodeExecutionCapabilities := &risk.Event{
		ID:          ThreatActorEstablishesCodeExecutionCapabilitiesID,
		Name:        "Threat Actor Establishes Code Execution Capabilities",
		Description: "A threat actor establishes code execution capabilities on a system.",
		Probability: &risk.Probability{
			// The Threat Detection & response team has never observed a threat actor establishing code execution capabilities and is pretty confident in their host defenses
			// However, they're not sure what the threat actor is capable of or what would really happen if they did establish code execution capabilities, so they're not very confident in their predictions
			// The team assessed this prediction on a yearly basis, as they don't expect to see this event very often
			ExpectedFrequency: "yearly",
			Minimum:           0.1,
			MinimumConfidence: 0.5,
			Maximum:           0.5,
			MaximumConfidence: 0.5,
		},
		Impact: []*risk.Impact{
			{
				// The Threat Detection & Response group wants to estimate the number of threat actors with access to the network as a result of establishing code execution capabilities.
				// The team knows that threat actors often sell access to networks, the team estimates between 1 and 5 unique threat actors will have access to the network but they're not very confident in their predictions
				// The team also knows that the threat actor cal sell access to the network multiple times, so they've marked their minimum and maximum impact events at 1 and 5 respectively
				ImpactID:                              impactThreatActorsWithAccessID,
				Name:                                  "Threat Actors with Access",
				Unit:                                  "Threat Actor",
				PositiveImpact:                        false,
				Description:                           "The number of threat actors with access to the network as a result of establishing code execution capabilities.",
				ExpectedFrequency:                     "yearly",
				MinimumIndividualUnitImpact:           1,
				MinimumIndividualUnitImpactConfidence: 0.9,
				MaximumIndividualUnitImpact:           5,
				MaximumIndividualUnitImpactConfidence: 0.5,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.9,
				MaximumImpactEvents:                   5,
				MaximumImpactEventsConfidence:         0.5,
			},
		},
		Dependencies: []*risk.Dependency{
			{
				// The employee has to accept a malicious Duo Push notification for the threat actor to establish code execution capabilities
				DependsOnEventID: EmployeeAcceptsMaliciousDuoPushID,
				Happens:          true,
			},
			{
				// the host-based controls have to not catch the malicious activity or code for the threat actor to establish code execution capabilities
				DependsOnEventID: HostBasedControlsCatchMaliciousActivityOrCodeID,
				Happens:          false,
			},
			{
				// The behavioral controls have to not catch the anomalous account behavior for the threat actor to establish code execution capabilities
				DependsOnEventID: BehavioralControlsCatchAnomalousAccountBehaviorID,
				Happens:          false,
			},
		},
	}

	ThreatActorDeliversRansomwarePayload := &risk.Event{
		ID:          ThreatActorDeliversRansomwarePayloadID,
		Name:        "Threat Actor Delivers Ransomware Payload",
		Description: "A threat actor delivers a ransomware payload to a system.",
		Probability: &risk.Probability{
			// The Threat Detection and Response group is pretty sure that at between 60% and 90% of the time, when possible a threat group will deliver ransomware to a system
			// The team is confident in their predictions, so they've marked both their minimum and maximum predictions at 80% confidence
			// The team assessed this prediction on a yearly basis, as they don't expect to see this event very often
			ExpectedFrequency: "yearly",
			Minimum:           0.6,
			MinimumConfidence: 0.8,
			Maximum:           0.9,
			MaximumConfidence: 0.8,
		},
		Impact: []*risk.Impact{
			{
				// The Threat Detection & Response group wants to estimate the number of malware instances on the network as a result of a ransomware payload being delivered.
				// The team knows that threat actors may deliver ransomware to multiple systems at once or multiple different ransomware payloads to the same system
				// since the team believes on average users have access to between 1 and 5 systems, they've marked their minimum and maximum impact events at 1 and 5 respectively
				// The team has read that on average threat actors will drop between 1 and 3 different ransomware payloads on a system, so they've marked their minimum and maximum individual unit impacts at 1 and 3 respectively
				ImpactID:                              impactMalwareOnNetworkID,
				Name:                                  "Malware on Network",
				Unit:                                  "Malware Instance",
				PositiveImpact:                        false,
				Description:                           "The number of malware instances on the network as a result of a ransomware payload being delivered.",
				ExpectedFrequency:                     "yearly",
				MinimumIndividualUnitImpact:           1,
				MinimumIndividualUnitImpactConfidence: 0.9,
				MaximumIndividualUnitImpact:           3,
				MaximumIndividualUnitImpactConfidence: 0.6,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.9,
				MaximumImpactEvents:                   5,
				MaximumImpactEventsConfidence:         0.9,
			},
		},
		Dependencies: []*risk.Dependency{
			{
				// The threat actor has to establish code execution capabilities for them to deliver a ransomware payload
				DependsOnEventID: ThreatActorEstablishesCodeExecutionCapabilitiesID,
				Happens:          true,
			},
			{
				// The network-based controls have to not catch the malicious command and control traffic for the threat actor to deliver a ransomware payload
				DependsOnEventID: NetworkBasedControlsCatchMaliciousCommandAndControlTrafficID,
				Happens:          false,
			},
			{
				// The host-based controls have to not catch the malicious activity or code for the threat actor to deliver a ransomware payload
				DependsOnEventID: HostBasedControlsCatchMaliciousActivityOrCodeID,
				Happens:          false,
			},
			{
				// The behavioral controls have to not catch the anomalous account behavior for the threat actor to deliver a ransomware payload
				DependsOnEventID: BehavioralControlsCatchAnomalousAccountBehaviorID,
				Happens:          false,
			},
		},
	}

	RansomwarePropogatesThroughoutNetworkWithoutDetection := &risk.Event{
		ID:          RansomwarePropogatesThroughoutNetworkWithoutDetectionID,
		Name:        "Ransomware Propogates Throughout Network Without Detection",
		Description: "Ransomware propogates throughout the network without detection.",
		Probability: &risk.Probability{
			// The Threat Detection and Response group is pretty sure that at between 10% and 40% of the time, ransomware will propogate throughout the network without detection
			// Overall they're pretty confident in their controls but they have some uncertainty and known gaps so they're not very confident in their predictions
			// The team assessed this prediction on a yearly basis, as they don't expect to see this event very often
			ExpectedFrequency: "yearly",
			Minimum:           0.1,
			MinimumConfidence: 0.6,
			Maximum:           0.4,
			MaximumConfidence: 0.6,
		},
		Impact: []*risk.Impact{
			{
				// The Threat Detection & Response group wants to estimate the number of lateral movement events as a result of ransomware propogating throughout the network without detection.
				// For any given event, the team believes ransomware can gain access to anywhere from 1 to 100 systems depending on what sort of access they identify
				// this is because the organization knows developers often keep keys, etc. on their systems or application servers
				// The team is very confident that at least 1 lateral movement event would occur, but they think theres not much chance of 100 systems being hit at once so they've marked their confidences accordingly
				// The team also thinks that ransomware can move laterally to multiple systems at once, so they've marked their minimum and maximum impact events at 1 and 100 respectively
				ImpactID:                              impactLateralMovementEventsID,
				Name:                                  "Lateral Movement Events",
				Unit:                                  "Lateral Movement Event",
				PositiveImpact:                        false,
				Description:                           "The number of lateral movement events as a result of ransomware propogating throughout the network without detection.",
				ExpectedFrequency:                     "yearly",
				MinimumIndividualUnitImpact:           1,
				MinimumIndividualUnitImpactConfidence: 0.9,
				MaximumIndividualUnitImpact:           100,
				MaximumIndividualUnitImpactConfidence: 0.5,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.9,
				MaximumImpactEvents:                   100,
				MaximumImpactEventsConfidence:         0.5,
			},
		},
		Dependencies: []*risk.Dependency{
			{
				// The threat actor has to deliver a ransomware payload for the ransomware to propogate throughout the network without detection
				DependsOnEventID: ThreatActorDeliversRansomwarePayloadID,
				Happens:          true,
			},
			{
				// The network-based controls have to not catch the malicious command and control traffic for the ransomware to propogate throughout the network without detection
				DependsOnEventID: NetworkBasedControlsCatchMaliciousCommandAndControlTrafficID,
				Happens:          false,
			},
			{
				// The host-based controls have to not catch the malicious activity or code for the ransomware to propogate throughout the network without detection
				DependsOnEventID: HostBasedControlsCatchMaliciousActivityOrCodeID,
				Happens:          false,
			},
			{
				// The behavioral controls have to not catch the anomalous account behavior for the ransomware to propogate throughout the network without detection
				DependsOnEventID: BehavioralControlsCatchAnomalousAccountBehaviorID,
				Happens:          false,
			},
		},
	}

	// Top Level Event
	MajorRansomwareEvent := &risk.Event{
		ID:          MajorRansomwareEventID,
		Name:        "Major Ransomware Event",
		Description: "A major ransomware event occurs.",
		Probability: &risk.Probability{
			// The team is pretty sure a major ransomware outbreak would be detected most of the time before it gets our of hand, but they're not 100% sure
			// The team has low confidence in their predictions but they're pretty sure that a major ransomware event would be detected at least 70% of the time, and at most 90% of the time
			// The team assessed this prediction on a yearly basis, as they don't expect to see this event very often
			ExpectedFrequency: "yearly",
			Minimum:           0.1,
			MinimumConfidence: 0.1,
			Maximum:           0.3,
			MaximumConfidence: 0.1,
		},
		Impact: []*risk.Impact{
			{
				// The organization wants to estimate the cost of rebuilding the network as a result of a major ransomware event.
				// The team knows that rebuilding the network will cost between $1,000,000 and $10,000,000 based on the size of the organization and the amount of data lost
				// The team is confident in their predictions
				ImpactID:                              impactRebuildingNetworkID,
				Name:                                  "Rebuilding Network",
				Unit:                                  "USD",
				PositiveImpact:                        false,
				Description:                           "The cost of rebuilding the network as a result of a major ransomware event.",
				ExpectedFrequency:                     "yearly",
				MinimumIndividualUnitImpact:           1000000,
				MinimumIndividualUnitImpactConfidence: 0.9,
				MaximumIndividualUnitImpact:           10000000,
				MaximumIndividualUnitImpactConfidence: 0.9,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.9,
				MaximumImpactEvents:                   1,
				MaximumImpactEventsConfidence:         0.9,
			},
			{
				// The organization wants to estimate the cost of customer renumeration as a result of a major ransomware event.
				// The organization estimates that between 1 and 25000 customers would need to be renumerated based on the amount of data lost
				// The organization estimates renumeration costs per customer at between $100 and $1000
				// The team is fairly confident in their predictions
				ImpactID:                              impactCustomerRenumerationID,
				Name:                                  "Customer Renumeration",
				Unit:                                  "USD",
				PositiveImpact:                        false,
				Description:                           "The cost of customer renumeration as a result of a major ransomware event.",
				ExpectedFrequency:                     "yearly",
				MinimumIndividualUnitImpact:           100,
				MinimumIndividualUnitImpactConfidence: 0.7,
				MaximumIndividualUnitImpact:           1000,
				MaximumIndividualUnitImpactConfidence: 0.7,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.7,
				MaximumImpactEvents:                   25000,
				MaximumImpactEventsConfidence:         0.7,
			},
			{
				// The organization wants to estimate the cost of customer loss as a result of a major ransomware event.
				// The organization estimates that between 1 and 10000 customers would stop using the service based on the amount of data lost
				// Each loss would cost the organization between $10 and $100 depending on the customer's subscription level
				// subscriptions are paid on a monthly basis
				// The team is fairly confident in their customer predictions but knows the subscription costs of their customers are set by the organization so they are very confident in their subscription cost predictions
				ImpactID:                              impactCustomerLossID,
				Name:                                  "Customer Loss",
				Unit:                                  "USD",
				PositiveImpact:                        false,
				Description:                           "The cost of customer loss as a result of a major ransomware event.",
				ExpectedFrequency:                     "monthly",
				MinimumIndividualUnitImpact:           10,
				MinimumIndividualUnitImpactConfidence: 1,
				MaximumIndividualUnitImpact:           100,
				MaximumIndividualUnitImpactConfidence: 1,
				MinimumImpactEvents:                   1,
				MinimumImpactEventsConfidence:         0.7,
				MaximumImpactEvents:                   10000,
				MaximumImpactEventsConfidence:         0.7,
			},
		},
		Dependencies: []*risk.Dependency{
			{
				// The ransomware has to propogate throughout the network without detection for a major ransomware event to occur
				DependsOnEventID: RansomwarePropogatesThroughoutNetworkWithoutDetectionID,
				Happens:          true,
			},
		},
	}

	// --- RISK ASSESSMENT ---

	Events := []*risk.Event{
		AntiPhishFilter,
		EmployeeReportsPhish,
		BehavioralControlsCatchAnomalousAccountBehavior,
		HostBasedControlsCatchMaliciousActivityOrCode,
		NetworkBasedControlsCatchMaliciousCommandAndControlTraffic,
		EmployeeGetsPhished,
		EmployeeAcceptsMaliciousDuoPush,
		AttemptedPhishingEmail,
		ThreatActorEstablishesCodeExecutionCapabilities,
		ThreatActorDeliversRansomwarePayload,
		RansomwarePropogatesThroughoutNetworkWithoutDetection,
		MajorRansomwareEvent,
	}

	ProbabilityMap, ImpactMap, err := analysis.MonteCarlo(Events, 100_000)

	if err != nil {
		panic(fmt.Errorf("error running monte carlo analysis: %w", err))
	}

	ProbMap := make(map[string]float64)
	for k, v := range ProbabilityMap {
		for _, e := range Events {
			if e.ID == k {
				ProbMap[e.Name] = v
			}
		}
	}

	jProbabilityMap, err := json.Marshal(ProbMap)
	if err != nil {
		panic(err)
	}

	jImpactMap, err := json.Marshal(ImpactMap)
	if err != nil {
		panic(err)
	}

	pmf_file, err := os.Create("probabilities.json")
	if err != nil {
		panic(err)
	}

	_, err = pmf_file.Write(jProbabilityMap)
	if err != nil {
		panic(err)
	}

	impact_file, err := os.Create("impacts.json")
	if err != nil {
		panic(err)
	}

	_, err = impact_file.Write(jImpactMap)
	if err != nil {
		panic(err)
	}

}
