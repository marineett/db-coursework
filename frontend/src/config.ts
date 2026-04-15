export const API_ENDPOINTS = {
    AUTH: {
        LOGIN: '/api/auth/authorize',
        REGISTER: {
            CLIENT: '/api/registration/client',
            REPETITOR: '/api/registration/repetitor',
            MODERATOR: '/api/registration/moderator',
            ADMIN: '/api/registration/admin',
        }
    },
    CLIENT: {
        CONTRACTS: (clientId: number, offset: number, limit: number, status: number) => 
            `/api/client/get_contracts?client_id=${clientId}&offset=${offset}&limit=${limit}&status=${status}`,
        CREATE_CONTRACT: '/api/client/create_contract',
        MAKE_REVIEW: '/api/client/make_review',
        GET_REVIEWS: (repetitorId: number) => `/api/contract/get_reviews?repetitor_id=${repetitorId}`,
        GET_REVIEW: '/api/contract/get_review'
    },
    REPETITOR: {
        GET_CONTRACTS: (repetitorId: number, offset: number, limit: number, status: number) =>
            `/api/repetitor/get_contracts?repetitor_id=${repetitorId}&offset=${offset}&limit=${limit}&status=${status}`,
        GET_AVAILABLE_CONTRACTS: (offset: number, limit: number, status: number) =>
            `/api/repetitor/get_available_contracts?offset=${offset}&limit=${limit}&status=${status}`,
        RESPOND_TO_CONTRACT: (contractId: number) =>
            `/api/repetitor/respond_to_contract?contract_id=${contractId}`,
        ACCEPT_CONTRACT: (repetitorId: number, contractId: number) =>
            `/api/repetitor/accept_contract?repetitor_id=${repetitorId}&contract_id=${contractId}`,
        GET_REVIEWS: (repetitorId: number) => `/api/repetitor/get_reviews?repetitor_id=${repetitorId}`,
        CONTRACTS: (repetitorId: number, offset: number, limit: number, status: number) => 
            `/api/contract/get_contracts?repetitor_id=${repetitorId}&offset=${offset}&limit=${limit}&status=${status}`,
        UPDATE_CONTRACT_STATUS: '/api/contract/update_status',
        MAKE_REVIEW: '/api/repetitor/make_review',
        PAY_FOR_CONTRACT: (contractId: number, amount: number) =>
            `/api/repetitor/pay_for_contract?contract_id=${contractId}&amount=${amount}`
    },
    MODERATOR: {
        CONTRACTS: (moderatorId: number, offset: number, limit: number) => 
            `/api/moderator/contracts?moderator_id=${moderatorId}&offset=${offset}&limit=${limit}`,
        GET_TRANSACTION_TO_APPROVE: '/api/moderator/get_transaction_to_approve',
        APPROVE_TRANSACTION: (transactionId: number) =>
            `/api/moderator/approve_transaction?transaction_id=${transactionId}`,
        REJECT_TRANSACTION: (transactionId: number) =>
            `/api/moderator/reject_transaction?transaction_id=${transactionId}`
    },
    ADMIN: {
        CONTRACTS: (adminId: number, offset: number, limit: number) => 
            `/api/admin/contracts?admin_id=${adminId}&offset=${offset}&limit=${limit}`,
    },
    CHAT: {
        START_CR_CHAT: (clientId: number, repetitorId: number) =>
            `/api/chat/start_cr_chat?c_id=${clientId}&r_id=${repetitorId}`,
        GET_CHAT: (chatId: number) =>
            `/api/chat/get_chat?id=${chatId}`,
        SEND_MESSAGE: (senderId: number, chatId: number) =>
            `/api/chat/send_message?sender_id=${senderId}&chat_id=${chatId}`,
        GET_MESSAGES: (chatId: number, offset: number, limit: number) =>
            `/api/chat/get_messages?id=${chatId}&messages_offset=${offset}&messages_limit=${limit}`
    }
}; 