# demo
# running 5 pods on 3 nodes, 4GB 1CPU each


# find which node is the leader
for i in 0 1 2 3 4; do
    echo "=== zk-$i ==="
    kubectl exec zk-$i -- zkServer.sh status | grep Mode
done
