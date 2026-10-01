#include <algorithm>
#include <functional>
#include <string>
#include <unordered_map>
#include <vector>
using namespace std;

vector<int> keepers(vector<vector<string>>& contacts) {
    int n = contacts.size();
    // Contacts and addresses are nodes; search from each unseen contact in
    // number order, so the first contact of a person is its lowest.
    unordered_map<string, vector<int>> withAddress;
    for (int i = 0; i < n; i++)
        for (auto& a : contacts[i]) withAddress[a].push_back(i);
    vector<int> out(n, -1);
    for (int start = 0; start < n; start++) {
        if (out[start] >= 0) continue;
        vector<int> stack{start};
        out[start] = start;
        while (!stack.empty()) {
            int cur = stack.back();
            stack.pop_back();
            for (auto& a : contacts[cur])
                for (int other : withAddress[a])
                    if (out[other] < 0) {
                        out[other] = start;
                        stack.push_back(other);
                    }
        }
    }
    return out;
}
