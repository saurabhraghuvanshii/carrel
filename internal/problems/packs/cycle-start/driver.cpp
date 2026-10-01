#include <algorithm>
#include <climits>
#include <cmath>
#include <cstdlib>
#include <exception>
#include <iomanip>
#include <iostream>
#include <list>
#include <map>
#include <queue>
#include <set>
#include <sstream>
#include <stdexcept>
#include <stack>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>
using namespace std;

struct ListNode {
    int val;
    ListNode* next;
    ListNode(int v) : val(v), next(nullptr) {}
};

#include "solution.cpp"

vector<int> readArray() {
    int n;
    cin >> n;
    vector<int> a(n);
    for (auto& x : a) cin >> x;
    return a;
}

string show(const vector<int>& a) {
    ostringstream o;
    o << '[';
    for (size_t i = 0; i < a.size(); i++) {
        if (i > 0) o << ", ";
        o << a[i];
    }
    o << ']';
    return o.str();
}

// Reads T cases from stdin and prints one line per case. Each case is a list, then the position the last node points back to, or -1.
// Flushes after every case so a crash cannot hide finished results.
// Anything the learner prints with cout goes to stderr.
int main() {
    streambuf* real = cout.rdbuf();
    cout.rdbuf(cerr.rdbuf());
    ostream out(real);
    int t;
    cin >> t;
    for (int c = 0; c < t; c++) {
        vector<int> vals = readArray();
        int pos;
        cin >> pos;
        vector<ListNode*> nodes;
        for (int v : vals) {
            nodes.push_back(new ListNode(v));
            if (nodes.size() > 1) nodes[nodes.size() - 2]->next = nodes.back();
        }
        if (pos >= 0) nodes.back()->next = nodes[pos];
        ListNode* head = nodes.empty() ? nullptr : nodes[0];
        try {
            ListNode* start = loopStart(head);
            int at = -1;
            for (size_t i = 0; i < nodes.size(); i++) {
                if (nodes[i] == start) at = i;
            }
            if (start != nullptr && at < 0) {
                out << "ERROR returned a node that is not in the list\n";
            } else {
                out << at << '\n';
            }
        } catch (const exception& e) {
            out << "ERROR " << e.what() << '\n';
        }
        out.flush();
    }
    return 0;
}
