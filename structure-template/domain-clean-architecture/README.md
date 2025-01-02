# domain 層：應該只包含核心業務模型和行為的定義，例如 Entity、Value Object 和 interface（如 Repository 和 Service interface）。
# application 層：應該負責調用 domain 層的接口，協調應用的操作，並實現具體的業務邏輯。