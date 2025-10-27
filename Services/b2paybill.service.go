package Services

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	abstracts "github.com/venomous-maker/go-mpesa/Abstracts"
)

// BusinessToPayBillService handles Business-to-PayBill (B2B PayBill) payments.
// This service allows a business to pay directly to a PayBill number or store on behalf of a consumer.
type BusinessToPayBillService struct {
	Config                  *abstracts.MpesaConfig
	Client                  abstracts.MpesaInterface
	initiator               string
	commandID               string
	senderIdentifierType    string
	recipientIdentifierType string
	amount                  float64
	partyA                  string
	partyB                  string
	accountReference        string
	requester               string
	remarks                 string
	occasion                string
	response                map[string]any
}

// NewBusinessToPayBillService creates a new B2B PayBill service instance.
func NewBusinessToPayBillService(cfg *abstracts.MpesaConfig, client abstracts.MpesaInterface) *BusinessToPayBillService {
	return &BusinessToPayBillService{
		Config:                  cfg,
		Client:                  client,
		commandID:               "BusinessPayBill",
		senderIdentifierType:    "4",
		recipientIdentifierType: "4",
	}
}

// SetInitiator sets the initiator (operator username) for the transaction.
func (s *BusinessToPayBillService) SetInitiator(name string) *BusinessToPayBillService {
	s.initiator = name
	return s
}

// SetSecurityCredential encrypts and sets the security credential (initiator password).
func (s *BusinessToPayBillService) SetSecurityCredential(password string) error {
	return s.Config.SetSecurityCredential(password)
}

// SetAmount sets the transaction amount in KES.
func (s *BusinessToPayBillService) SetAmount(amount float64) *BusinessToPayBillService {
	s.amount = amount
	return s
}

// SetPartyA sets the shortcode from which money will be deducted.
func (s *BusinessToPayBillService) SetPartyA(code string) *BusinessToPayBillService {
	s.partyA = code
	s.Config.SetBusinessCode(code)
	return s
}

// SetPartyB sets the destination shortcode (paybill) to which money will be credited.
func (s *BusinessToPayBillService) SetPartyB(code string) *BusinessToPayBillService {
	s.partyB = code
	return s
}

// SetAccountReference sets an account/reference associated with the payment.
func (s *BusinessToPayBillService) SetAccountReference(ref string) *BusinessToPayBillService {
	s.accountReference = ref
	return s
}

// SetRequester sets the optional consumer mobile number on whose behalf the payment is made.
func (s *BusinessToPayBillService) SetRequester(msisdn string) *BusinessToPayBillService {
	s.requester = msisdn
	return s
}

// SetRemarks sets transaction remarks.
func (s *BusinessToPayBillService) SetRemarks(r string) *BusinessToPayBillService {
	s.remarks = r
	return s
}

// SetOccasion sets the occasion for the transaction.
func (s *BusinessToPayBillService) SetOccasion(o string) *BusinessToPayBillService {
	s.occasion = o
	return s
}

// SetQueueTimeoutURL updates the config queue timeout URL.
func (s *BusinessToPayBillService) SetQueueTimeoutURL(url string) *BusinessToPayBillService {
	s.Config.SetQueueTimeoutURL(url)
	return s
}

// SetResultURL updates the config result URL.
func (s *BusinessToPayBillService) SetResultURL(url string) *BusinessToPayBillService {
	s.Config.SetResultURL(url)
	return s
}

// Send constructs and sends the B2B BusinessPayBill payment request to M-Pesa.
func (s *BusinessToPayBillService) Send() (map[string]any, error) {
	// Validate required fields
	if s.initiator == "" {
		return nil, errors.New("initiator is required")
	}
	if s.Config.GetSecurityCredential() == "" {
		return nil, errors.New("security credential is required; call SetSecurityCredential")
	}
	if s.amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if s.partyA == "" && s.Config.GetBusinessCode() == "" {
		return nil, errors.New("partyA (business shortcode) is required")
	}
	if s.partyB == "" {
		return nil, errors.New("partyB (destination shortcode/paybill) is required")
	}

	data := map[string]any{
		"Initiator":              s.initiator,
		"SecurityCredential":     s.Config.GetSecurityCredential(),
		"CommandID":              s.commandID,
		"SenderIdentifierType":   s.senderIdentifierType,
		"RecieverIdentifierType": s.recipientIdentifierType,
		"Amount":                 math.Round(s.amount),
		"PartyA":                 s.getPartyA(),
		"PartyB":                 s.partyB,
		"AccountReference":       s.accountReference,
		"Requester":              s.requester,
		"Remarks":                s.remarks,
		"QueueTimeOutURL":        s.Config.GetQueueTimeoutURL(),
		"ResultURL":              s.Config.GetResultURL(),
		"Occasion":               s.occasion,
	}

	resp, err := s.Client.ExecuteRequest(data, "/mpesa/b2b/v1/paymentrequest")
	if err != nil {
		return nil, err
	}

	s.response = resp
	return resp, nil
}

func (s *BusinessToPayBillService) getPartyA() string {
	if s.partyA != "" {
		return s.partyA
	}
	return s.Config.GetBusinessCode()
}

// GetResponse returns the last API response stored by the service.
func (s *BusinessToPayBillService) GetResponse() map[string]any {
	return s.response
}

// B2PayBillCallbackResult represents a parsed B2PayBill callback payload.
// It normalizes common fields and provides easy access to ResultParameters and ReferenceData.
type B2PayBillCallbackResult struct {
	ResultCode               string // numeric result code as string
	ResultDesc               string // human readable description
	TransactionID            string // M-Pesa transaction ID
	OriginatorConversationID string
	ConversationID           string
	ResultParameters         map[string]string // key->value map from ResultParameters.ResultParameter
	ReferenceData            map[string]string // key->value map from ReferenceData.ReferenceItem
	Raw                      map[string]any    // original payload
	Success                  bool              // true if ResultCode == 0
}

// ParseCallback parses a B2PayBill callback payload and returns a structured result.
// The method tolerates variations in payload shapes (slices vs single objects) and
// attempts to extract ResultCode, ResultDesc, TransactionID, ResultParameters and ReferenceData.
func (s *BusinessToPayBillService) ParseCallback(payload map[string]any) (*B2PayBillCallbackResult, error) {
	res := &B2PayBillCallbackResult{
		ResultParameters: make(map[string]string),
		ReferenceData:    make(map[string]string),
		Raw:              payload,
	}

	// locate "Result" node (case-insensitive tolerant)
	var resultNode any
	if v, ok := payload["Result"]; ok {
		resultNode = v
	} else if v, ok := payload["result"]; ok {
		resultNode = v
	} else {
		return nil, errors.New("payload missing Result node")
	}

	resultMap, ok := resultNode.(map[string]any)
	if !ok {
		// try map[string]interface{} conversion
		if m2, ok2 := resultNode.(map[string]interface{}); ok2 {
			resultMap = make(map[string]any)
			for k, v := range m2 {
				resultMap[k] = v
			}
		} else {
			return nil, errors.New("Result node is not an object")
		}
	}

	res.ResultCode = toString(resultMap["ResultCode"]) // may be string or number
	res.ResultDesc = toString(resultMap["ResultDesc"])
	res.TransactionID = toString(resultMap["TransactionID"])
	res.OriginatorConversationID = toString(resultMap["OriginatorConversationID"])
	res.ConversationID = toString(resultMap["ConversationID"])

	// Parse ResultParameters -> ResultParameter (array or single)
	if rpRaw, ok := resultMap["ResultParameters"]; ok {
		switch rp := rpRaw.(type) {
		case map[string]any:
			parseResultParameterArray(rp["ResultParameter"], res.ResultParameters)
		default:
			parseResultParameterArray(rpRaw, res.ResultParameters)
		}
	}

	// Parse ReferenceData -> ReferenceItem
	if rdRaw, ok := resultMap["ReferenceData"]; ok {
		switch rd := rdRaw.(type) {
		case map[string]any:
			parseReferenceItem(rd["ReferenceItem"], res.ReferenceData)
		default:
			parseReferenceItem(rdRaw, res.ReferenceData)
		}
	}

	// determine success
	if i, err := strconv.Atoi(res.ResultCode); err == nil {
		res.Success = i == 0
	} else {
		res.Success = res.ResultCode == "0"
	}

	return res, nil
}

// helpers
func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		// convert float to int string if it's integer-valued
		if t == math.Trunc(t) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case nil:
		return ""
	default:
		return fmt.Sprint(t)
	}
}

func parseResultParameterArray(input any, out map[string]string) {
	if input == nil {
		return
	}

	// handle slice
	if arr, ok := input.([]any); ok {
		for _, item := range arr {
			if m, ok := item.(map[string]any); ok {
				k := toString(m["Key"])
				v := toString(m["Value"])
				if k != "" {
					out[k] = v
				}
			}
		}
		return
	}

	// handle []interface{}
	if arr2, ok := input.([]interface{}); ok {
		for _, item := range arr2 {
			if m, ok := item.(map[string]interface{}); ok {
				k := toString(m["Key"])
				v := toString(m["Value"])
				if k != "" {
					out[k] = v
				}
			}
		}
		return
	}

	// single object
	if m, ok := input.(map[string]any); ok {
		k := toString(m["Key"])
		v := toString(m["Value"])
		if k != "" {
			out[k] = v
		}
		return
	}
	if m2, ok := input.(map[string]interface{}); ok {
		k := toString(m2["Key"])
		v := toString(m2["Value"])
		if k != "" {
			out[k] = v
		}
		return
	}
}

func parseReferenceItem(input any, out map[string]string) {
	if input == nil {
		return
	}

	if arr, ok := input.([]any); ok {
		for _, item := range arr {
			if m, ok := item.(map[string]any); ok {
				k := toString(m["Key"])
				v := toString(m["Value"])
				if k != "" {
					out[k] = v
				}
			}
		}
		return
	}

	if arr2, ok := input.([]interface{}); ok {
		for _, item := range arr2 {
			if m, ok := item.(map[string]interface{}); ok {
				k := toString(m["Key"])
				v := toString(m["Value"])
				if k != "" {
					out[k] = v
				}
			}
		}
		return
	}

	if m, ok := input.(map[string]any); ok {
		k := toString(m["Key"])
		v := toString(m["Value"])
		if k != "" {
			out[k] = v
		}
		return
	}

	if m2, ok := input.(map[string]interface{}); ok {
		k := toString(m2["Key"])
		v := toString(m2["Value"])
		if k != "" {
			out[k] = v
		}
		return
	}
}
